package connection

import (
	"encoding/binary"
	"fmt"
	"github.com/joernweissenborn/aursir4go/aurarath"
	"github.com/joernweissenborn/stream2go"
	"github.com/pebbe/zmq4"
	"net"
	"sync"
	"errors"
	"time"
	"github.com/joernweissenborn/eventual2go"
)

type Outgoing struct {
	*sync.Mutex
	skt         *zmq4.Socket
	ipportbytes []byte
	targetAddress string
	targetPort int
}

func NewOutgoing(uuid string, targetAddress string, targetPort int) (out eventual2go.StreamController, err error) {
	var o Outgoing
	o.Mutex = new(sync.Mutex)

	o.skt, err = zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return
	}
	o.skt.SetIdentity(uuid)

	err = o.skt.Connect(fmt.Sprintf("tcp://%s:%d", targetAddress, targetPort))

	out = eventual2go.NewStreamController()
	out.Stream.Listen(o.send)
	out.Stream.Closed.Then(o.Close)
	return
}

func (o Outgoing) send(d eventual2go.Data) {
	o.Lock()
	defer o.Unlock()

	o.skt.SendMessage(d, 0)

	return
}
func (o Outgoing) Close(interface{}) interface{} {
//	log.Println("Stop")
	o.Lock()
	defer o.Unlock()
	return o.skt.Close()

}
