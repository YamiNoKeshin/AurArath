package connection

import (
	"github.com/pebbe/zmq4"
	"net"

	"fmt"
	"github.com/joernweissenborn/eventual2go"
	"time"
	"sync"
)

type Incoming struct {

	m *sync.Mutex

	addr string
	port uint16
	skt  *zmq4.Socket
	in   eventual2go.StreamController
	stopped bool
}

func NewIncoming(addr string) (i *Incoming, err error) {
	i = new(Incoming)
	i.m = new(sync.Mutex)
	i.addr = addr
	i.in = eventual2go.NewStreamController()
	err = i.setupSocket()
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
func (i *Incoming) setupSocket() (err error) {
	i.port = getRandomPort()
	i.skt, err = zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}
	err = i.skt.Bind(fmt.Sprintf("tcp://%s:%d",i.addr, i.port))
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

func (i *Incoming) listen() {
	poller := zmq4.NewPoller()
	poller.Add(i.skt, zmq4.POLLIN)

	for {
		i.m.Lock()
		if i.stopped {
			i.m.Unlock()

			return
		}
		sockets, err := poller.Poll(100*time.Millisecond)
		if err != nil {
			continue
		}
		for range sockets {

			msg, err := i.skt.RecvMessage(0)
			if err == nil {
				i.in.Add(Message{i.addr, msg})
			}
		}
		i.m.Unlock()
	}
}

func (i *Incoming) Close() {
	i.m.Lock()
	i.stopped = true
	defer i.m.Unlock()
	i.skt.Close()
}