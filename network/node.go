package network

import (
	"github.com/hashicorp/serf/command/agent"
	//"github.com/hashicorp/serf/serf"
	"io"
	//"net"
	//"os"
	//"github.com/mitchellh/cli"
)

type Node struct {
	agent           *agent.Agent
	shutdownChannel interface{}
}

func NewNode(interfaceName string, logger io.Writer) *Node {
	node := &Node{}
	node.shutdownChannel = make(chan interface{}, 1)
	/*c := &agent.Command{

		}*/
	return node
}
