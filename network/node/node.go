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
	"time"
	"sync"
)

type Node struct {
	mut *sync.RWMutex

	UUID string

	cfg *config.Config

	tags map[string]string

	agents     map[string]*agent.Agent
	mDNSAgents map[string] *agent.AgentMDNS

	eventHandler eventHandler

	knownPeers map[string][]string // A map from peer UUID to IPs
}

func New(cfg *config.Config, tags map[string]string) (node *Node) {
	node = new(Node)

	node.cfg = cfg
	node.tags = tags

	id, _ := uuid.NewV4()
	node.UUID = id.String()

	node.knownPeers = make(map[string][]string)
	node.mut = new(sync.RWMutex)
	node.eventHandler = newEventHandler()
	node.eventHandler.join.Listen(node.newPeer)
	node.eventHandler.leave.Listen(node.leftPeer)

	node.createSerfAgents()



	return node
}

func (n *Node) newPeer(d eventual2go.Data){
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

func (n *Node) Run()  {
	for _,agt := range n.agents {
		agt.Start()
	}
	n.createMDNSAgents()
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
	                                               // Set probe intervals that are aggressive for finding bad nodes
	serfConfig.MemberlistConfig.GossipInterval = 5 * time.Millisecond
	serfConfig.MemberlistConfig.ProbeInterval = 50 * time.Millisecond
	serfConfig.MemberlistConfig.ProbeTimeout = 25 * time.Millisecond
	serfConfig.MemberlistConfig.SuspicionMult = 1
	serfConfig.Init()
	agentConfig := agent.DefaultConfig()

	agt, err := agent.Create(agentConfig, serfConfig, n.cfg.Logger())

	if n.handleErr(err) {
		agt.RegisterEventHandler(n.eventHandler)
		n.agents[iface] = agt

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


func (n *Node) createMDNSAgents() {
	for _, iface := range n.cfg.NetworkInterfaces {
		n.createMDNSAgent(iface)
	}
}

func (n *Node) createMDNSAgent(ifaceAddr string) {

	ifaces,_ := net.Interfaces()
	index := -1
	for i, iface := range ifaces {
		addrs, _ := iface.Addrs()
		if len(addrs) != 0 {
			addr := strings.Split(addrs[0].String(),"/")[0]
			if addr == ifaceAddr {
				index = i
				break
			}
		}
	}

	iface:= ifaces[index]


	 var err error
	agt := n.agents[ifaceAddr]
	n.mDNSAgents[ifaceAddr], err = agent.NewAgentMDNS(agt, n.cfg.Logger(), false, n.UUID, "AurArath",
		&iface, net.ParseIP(agt.SerfConfig().MemberlistConfig.BindAddr), agt.SerfConfig().MemberlistConfig.BindPort)
	n.handleErr(err)

}

func (n *Node) QueryPeer(UUID string) {



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
