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

	iface, err := net.InterfaceByName(interfaceName)

	if err != nil {
		return nil
	}

	addresses, err := iface.Addrs()

	if err != nil {
		return nil
	}

	bind := addresses[0] // If the interface has more than one address, we should take the first one

	mDNSAgent, err := agent.NewAgentMDNS(serfAgent, logger, false, hostname, "AurArath", iface, net.ParseIP(bind.String()), 42000)

	if err != nil {
		return nil
	}

	node.agent = serfAgent
	node.mDNSAgent = mDNSAgent
	return &node

}
