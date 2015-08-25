package node

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"
	"net"
	"fmt"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/eventual2go"
	"strings"
	"sync"
	"github.com/joernweissenborn/aurarath/network/beacon"
	"time"
	"log"
	"io/ioutil"
)

type Node struct {
	mut *sync.RWMutex

	UUID string

	cfg *config.Config

	tags map[string]string

	agents     map[string]*agent.Agent

	beacons []*beacon.Beacon

	logger *log.Logger

	join eventual2go.StreamController
	leave eventual2go.StreamController
	query eventual2go.StreamController
}

func New(uuid string, cfg *config.Config, tags map[string]string) (node *Node) {

	node = new(Node)
	node.logger = log.New(cfg.Logger(),fmt.Sprintf("node %s ",uuid),log.Lshortfile)

	node.logger.Println("Initializing")

	node.cfg = cfg
	node.tags = tags

	node.UUID = uuid
	node.beacons = []*beacon.Beacon{}
	node.mut = new(sync.RWMutex)
	node.logger.Println("Launching Serf Agents")

	node.join = eventual2go.NewStreamController()
	node.join.First().Then(node.silenceBeacons)
	node.leave = eventual2go.NewStreamController()
	node.query = eventual2go.NewStreamController()

	node.createSerfAgents()


	return node
}


func (n *Node) silenceBeacons(eventual2go.Data) eventual2go.Data{
	for _,b := range n.beacons {
		b.Silence()
	}
	return nil
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
		n.logger.Println("Launching Agent on",iface)

		n.createSerfAgent(iface)
	}
}
func (n *Node) createSerfAgent(iface string) {

	serfConfig := serf.DefaultConfig()


	serfConfig.Tags = n.tags

	serfConfig.MemberlistConfig.BindAddr = iface
	serfConfig.MemberlistConfig.BindPort = getRandomPort(iface)

	serfConfig.NodeName = fmt.Sprintf("%s@%s",n.UUID,iface)
//	serfConfig.LogOutput = n.cfg.Logger()
	serfConfig.LogOutput = ioutil.Discard
//	serfConfig.MemberlistConfig.GossipInterval = 5 * time.Millisecond
//	serfConfig.MemberlistConfig.ProbeInterval = 50 * time.Millisecond
//	serfConfig.MemberlistConfig.ProbeTimeout = 25 * time.Millisecond
//	serfConfig.MemberlistConfig.SuspicionMult = 1
	serfConfig.Init()
	agentConfig := agent.DefaultConfig()
	agentConfig.Tags = n.tags
	agentConfig.LogLevel = "INFO"
	agt, err := agent.Create(agentConfig, serfConfig, ioutil.Discard)
//	agt, err := agent.Create(agentConfig, serfConfig, n.cfg.Logger())

	if n.handleErr(err) {
		eventHandler := newEventHandler()
		n.join.Join(eventHandler.join.WhereNot(n.isSelf))
		n.leave.Join(eventHandler.leave.WhereNot(n.isSelf).Transform(toLeaveEvent()))
		n.query.Join(eventHandler.query.Transform(toQueryEvent(iface)))
		agt.RegisterEventHandler(eventHandler)
		n.agents[iface] = agt
		n.logger.Println("Agent Created")
	} else {
		n.logger.Println("Failed to create Agent")
	}
}

func (n *Node) isSelf(d eventual2go.Data)bool{
	return strings.Contains(d.(serf.Member).Name,n.UUID)
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
//	n.logger.Printf("Found Peer at %s:%d",peeraddress.IP,peeraddress.Port)
	for _, agt := range n.agents {
		agt.Join([]string{fmt.Sprintf("%s:%d",peeraddress.IP,peeraddress.Port)},false)
	}
}

func (n *Node) Join() eventual2go.Stream {
	return n.join.Stream
}

func (n *Node) Leave() eventual2go.Stream {
	return n.leave.Stream
}

func (n *Node) Queries() eventual2go.Stream {
	return n.query.Stream
}

func (n *Node) Query(name string, data []byte, results eventual2go.StreamController) {
	wg := new(sync.WaitGroup)
	for iface, agt := range n.agents {
		params := &serf.QueryParam{FilterTags:n.tags,Timeout: 1*time.Second}
		resp, _ := agt.Query(name, data, params)
		wg.Add(1)
		go collectResponse(iface, resp, results, wg)
	}
	go waitForQueryFinish(results,wg)
}

func collectResponse(iface string,resp *serf.QueryResponse, s eventual2go.StreamController, wg *sync.WaitGroup){
	for r := range resp.ResponseCh() {
		s.Add(QueryResponseEvent{iface,r})
	}
	wg.Done()
}
func waitForQueryFinish(s eventual2go.StreamController, wg *sync.WaitGroup) {
	wg.Wait()
	s.Close()
}

func (n *Node) Shutdown() {
	for _, agt := range n.agents{
		n.logger.Println("shutting down")
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
