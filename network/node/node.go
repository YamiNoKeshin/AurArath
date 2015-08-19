package node

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"
	"net"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"

	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/eventual2go"
	"strings"
	"sync"
	"github.com/joernweissenborn/aurarath/network/beacon"
	"time"
	"log"
)

type Node struct {
	mut *sync.RWMutex

	UUID string

	cfg *config.Config

	tags map[string]string

	agents     map[string]*agent.Agent
	mDNSAgents map[string] *agent.AgentMDNS

	eventHandler eventHandler

	beacons []*beacon.Beacon

	knownPeers map[string][]string // A map from peer UUID to IPs

	logger *log.Logger
}

func New(cfg *config.Config, tags map[string]string) (node *Node) {
	node = new(Node)

	node.cfg = cfg
	node.tags = tags

	id, _ := uuid.NewV4()
	node.UUID = id.String()
	node.logger = log.New(cfg.Logger(),"Node",log.Lshortfile)
	node.knownPeers = make(map[string][]string)
	node.beacons = []*beacon.Beacon{}
	node.mut = new(sync.RWMutex)
	node.eventHandler = newEventHandler()
	node.eventHandler.join.Listen(node.newPeer)
	node.eventHandler.leave.Listen(node.leftPeer)

	node.createSerfAgents()



	return node
}

func (n *Node) newPeer(d eventual2go.Data){
	n.silenceBeacons()
	n.mut.Lock()
	defer  n.mut.Unlock()
	peerid := strings.Split(d.(serf.Member).Name,"@")[0]
	peerip :=d.(serf.Member).Addr.String()
	n.knownPeers[peerid] = append(n.knownPeers[peerid],peerip)
}
func (n *Node) leftPeer(d eventual2go.Data){
	n.mut.Lock()
	defer  n.mut.Unlock()
	peerid := strings.Split(d.(serf.Member).Name,"@")[0]
	peerip :=d.(serf.Member).Addr.String()
	old := n.knownPeers[peerid]
	n.knownPeers[peerid] = []string{}
	for _, ip := range old {
		if ip != peerip {
			n.knownPeers[peerid] = append(n.knownPeers[peerid],ip)
		}
	}
}

func (n *Node) silenceBeacons() {
	for _,b := range n.beacons {
		b.Silence()
	}
}

func (n *Node) Run()  {
	for _,agt := range n.agents {
		agt.Start()
	}
	n.launchBeacon()
	return
}

func (n *Node) createSerfAgents() {
	n.agents = make(map[string]*agent.Agent)
	for _, iface := range n.cfg.NetworkInterfaces {
		n.createSerfAgent(iface)
	}
}
func (n *Node) createSerfAgent(iface string) {

	n.mDNSAgents= make(map[string]*agent.AgentMDNS)

	serfConfig := serf.DefaultConfig()


	serfConfig.Tags = n.tags

	serfConfig.MemberlistConfig.BindAddr = iface
	serfConfig.MemberlistConfig.BindPort = getRandomPort(iface)

	serfConfig.NodeName = fmt.Sprintf("%s@%s",n.UUID,iface)
	serfConfig.LogOutput = n.cfg.Logger()
//	serfConfig.MemberlistConfig.GossipInterval = 5 * time.Millisecond
//	serfConfig.MemberlistConfig.ProbeInterval = 50 * time.Millisecond
//	serfConfig.MemberlistConfig.ProbeTimeout = 25 * time.Millisecond
//	serfConfig.MemberlistConfig.SuspicionMult = 1
	serfConfig.Init()
	agentConfig := agent.DefaultConfig()
	agentConfig.Tags = n.tags
	agt, err := agent.Create(agentConfig, serfConfig, n.cfg.Logger())

	if n.handleErr(err) {
		agt.RegisterEventHandler(n.eventHandler)
		n.agents[iface] = agt
	}
}

func (n *Node) launchBeacon() {
	n.logger.Println("Launching Beacon")
	listening := false
	for addr,agt := range n.agents {
		cfg := &beacon.Config{PingAddresses:[]string{addr},Port: 5557, PingInterval:500*time.Millisecond}
		port := uint16(agt.SerfConfig().MemberlistConfig.BindPort)
		n.logger.Println("Launching Beacon on",addr,port)
		payload := NewSignalPayload(port)
		n.logger.Println("Payload is",payload)
		b := beacon.New(payload,cfg)
		if !listening {
			n.logger.Println("Opening Listen",addr)
			b.Signals().Where(IsValidSignal).Transform(SignalToAdress).Listen(n.recvPeerSignal)
			b.Run()
			listening = true
		}
		b.Ping()
		n.beacons = append(n.beacons,b)
	}
}

func (n *Node) recvPeerSignal(d eventual2go.Data) {
	peeraddress := d.(PeerAddress)
	n.logger.Printf("Found Peer at %s:%d",peeraddress.IP,peeraddress.Port)
	for _, agt := range n.agents {
		agt.Join([]string{fmt.Sprintf("%s:%d",peeraddress.IP,peeraddress.Port)},false)
	}
}

func (n *Node) Join() eventual2go.Stream {
	return n.eventHandler.join.WhereNot(func(d eventual2go.Data)bool{
		return strings.Contains(d.(serf.Member).Name,n.UUID)
	})
}

func (n *Node) Leave() eventual2go.Stream {
	return n.eventHandler.leave
}

func (n *Node) Queries() eventual2go.Stream {
	return n.eventHandler.query
}

func (n *Node) Query(name string, data []byte, results eventual2go.StreamController) {
	wg := new(sync.WaitGroup)
	for _, agt := range n.agents {
		params := &serf.QueryParam{FilterTags:n.tags}
		resp, _ := agt.Query(name, data, params)
		wg.Add(1)
		go collectResponse(resp,results, wg)
	}
	go waitForQueryFinish(results,wg)
}

func collectResponse(resp *serf.QueryResponse, s eventual2go.StreamController, wg *sync.WaitGroup){
	for r := range resp.ResponseCh() {
		s.Add(r)
	}
	wg.Done()
}
func waitForQueryFinish(s eventual2go.StreamController, wg *sync.WaitGroup) {
	wg.Wait()
	s.Close()
}

func (n *Node) Shutdown() {
	for _, agt := range n.agents{
		agt.Leave()
		agt.Shutdown()
	}
}

func (n *Node) handleErr(err error) (ok bool){
	ok = err == nil
	if ok {return}
	n.cfg.Logger().Write([]byte(fmt.Sprintf("NodeInit Error %s",err)))
	panic(err)
	return
}

func getRandomPort(iface string) (int) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:0",iface))
	if err != nil {
		panic(fmt.Sprintf("Could not find free port on: %s",iface))
	}
	defer l.Close()
	return int(l.Addr().(*net.TCPAddr).Port)
}
