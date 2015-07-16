package network

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"
	"io"
	"net"
	"os"
)

type Node struct {
	agent     *agent.Agent
	mDNSAgent *agent.AgentMDNS
}

func NewNode(interfaceName string, logger io.Writer) *Node {
	node := Node{}
	serfConfig := serf.DefaultConfig()
	agentConfig := agent.DefaultConfig()

	serfAgent, err := agent.Create(agentConfig, serfConfig, logger)

	if err != nil {
		return nil
	}

	hostname, _ := os.Hostname()

	var iface *net.Interface
	iface, err = net.InterfaceByName(interfaceName)

	if err != nil {
		iface, err = net.InterfaceByIndex(0)

		if err != nil {
			//Log
		}
	}

	var addresses []net.Addr

	addresses, err = iface.Addrs()

	if err != nil {
		return nil
	}

	var bind net.Addr

	if len(addresses) > 1 {
		bind = addresses[0] // If the interface has more than one address, we should take the first one
	}

	var mDNSAgent agent.AgentMDNS

	mDNSAgent, err = agent.NewAgentMDNS(agent, logger, false, hostname, "AurArath", iface, net.ParseIP(bind.String()), 42000)

	if err != nil {
		return nil
	}

	node.agent = serfAgent
	node.mDNSAgent = mDNSAgent
	return &node

}
