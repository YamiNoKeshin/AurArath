package connection

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"sync"
	"github.com/joernweissenborn/eventual2go"
)

type Outgoing struct {
	*sync.Mutex
	skt         *zmq4.Socket
	ipportbytes []byte
	targetAddress string
	targetPort int
}

func NewOutgoing(uuid string, targetAddress string, targetPort int) (out *eventual2go.StreamController, err error) {

	var o Outgoing

	o.Mutex = new(sync.Mutex)

	o.skt, err = zmq4.NewSocket(zmq4.DEALER)

	if err != nil {
		return
	}

	err = o.skt.SetIdentity(uuid)
	if err != nil {
		return
	}

	err = o.skt.Connect(fmt.Sprintf("tcp://%s:%d", targetAddress, targetPort))

	if err != nil {
		return
	}
	out = eventual2go.NewStreamController()
	out.Stream.Listen(o.send)
	out.Stream.Closed().Then(o.Close)
	return
}

func (o Outgoing) send(d eventual2go.Data) {
	o.Lock()
	defer o.Unlock()

	_, err := o.skt.SendMessage(d)

	if err != nil {
		panic(err)
	}

	return
}
func (o Outgoing) Close(eventual2go.Data) eventual2go.Data {
	o.Lock()
	defer o.Unlock()
	return o.skt.Close()

}
