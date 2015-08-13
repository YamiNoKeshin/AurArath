package node

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"
	"net"
	"os"
	"aurarath/config"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"

)

type Node struct {
	cfg *config.Config

	tags []string

	agents     map[string]*agent.Agent
	mDNSAgents map[string] *agent.AgentMDNS

	UUID string
}

func New(cfg *config.Config, tags []string) (node *Node) {
	node = new(Node)

	node.cfg = cfg
	node.tags = tags

	node.UUID, _ = uuid.NewV4()

	node.createSerfAgents()

	node.createMDNSAgents()

	return &node
}

func (n *Node) createSerfAgents() {
	for _, iface := range n.cfg.NetworkInterfaces {
		n.createSerfAgent(iface)
	}
}
func (n *Node) createSerfAgent(iface string) {

	serfConfig := serf.DefaultConfig()

	serfConfig.Init()

	serfConfig.Tags = n.tags

	serfConfig.MemberlistConfig.BindAddr = iface
	var err error
	serfConfig.MemberlistConfig.BindPort, err = getRandomPort(iface)
	if !n.handleErr(err) {
		return
	}
	serfConfig.NodeName = n.UUID

	agentConfig := agent.DefaultConfig()

	agent, err := agent.Create(agentConfig, serfConfig, n.cfg.Logger())
	if n.handleErr(err) {
		n.agents[iface] = agent
	}
}


func (n *Node) createMDNSAgents() {
	for _, iface := range n.cfg.NetworkInterfaces {
		n.createMDNSAgent(iface)
	}
}

func (n *Node) createMDNSAgent(ifaceAddr string) {
	iface, err := net.InterfaceByName(ifaceAddr)
	if !n.handleErr(err) {
		return
	}

	hostname, err := os.Hostname()
	if !n.handleErr(err) {
		return
	}
	agt := n.agents[ifaceAddr]
	n.mDNSAgents[ifaceAddr], err = agent.NewAgentMDNS(agt, n.cfg.Logger(), false, hostname, "AurArath", iface, net.ParseIP(agt.SerfConfig().MemberlistConfig.BindAddr), agt.SerfConfig().MemberlistConfig.BindPort)
}

func (n *Node) handleErr(err error) (ok bool){
	ok = err == nil
	if ok {return}
	n.cfg.Logger().Write([]byte{fmt.Sprintf("NodeInit Error %s",err)})
}


func getRandomPort(iface string) (int, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:0",iface)) // listen on localhost
	defer l.Close()
	return int(l.Addr().(*net.TCPAddr).Port), err
}
