package connection

import (
	"github.com/pebbe/zmq4"
	"net"

	"fmt"
	"github.com/joernweissenborn/eventual2go"
)

type Incoming struct {
	port uint16
	skt  *zmq4.Socket
	in   eventual2go.StreamController
}

func NewIncoming(addr string) (i Incoming, err error) {
	i.in = eventual2go.NewStreamController()
	err = i.setupSocket(addr)
	if err == nil {
		go i.listen()
	}
	return
}

func (i Incoming) In() eventual2go.Stream {
	return i.in.Stream
}

func (i Incoming) Port() (port uint16) {
	return i.port
}
func (i *Incoming) setupSocket(addr string) (err error) {
	i.port = getRandomPort()
	i.skt, err = zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}
	err = i.skt.Bind(fmt.Sprintf("tcp://%s:%d",addr, i.port))
	return
}

func getRandomPort() uint16 {
	l, err := net.Listen("tcp", ":0") // listen on address
	if err != nil {
		panic(fmt.Sprintf("Could not find a free port %v",err))
	}
	defer l.Close()
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func (i Incoming) listen() {

	defer i.skt.Close()

	for {
		msg, err := i.skt.RecvMessage(0)
		if err == nil {
			i.in.Add(msg)
		}
	}

}
