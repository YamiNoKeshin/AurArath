package connection

import (
	"github.com/pebbe/zmq4"
	"net"

	"fmt"
	"github.com/joernweissenborn/stream2go"
)

type Incoming struct {
	port uint16
	skt  *zmq4.Socket
	in   stream2go.StreamController
}

func NewIncoming() (i Incoming, err error) {
	i.in = stream2go.New()
	err = i.setupSocket()
	if err != nil {
		go i.Listen()
	}
	return
}
func (i *Incoming) setupSocket() (err error) {
	i.port = getRandomPort()
	i.skt, err = zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}
	err = i.skt.Bind(fmt.Sprintf("tcp://*:%d", i.port))
	return
}

func getRandomPort() uint16 {
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
		panic("Could not find a free port")
	}
	defer l.Close()
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func (i Incoming) Listen() {

	defer i.skt.Close()

	for {
		msg, err := i.skt.RecvMessage(0)
		if err == nil {
			i.in.Add(msg)
		}
	}

}
