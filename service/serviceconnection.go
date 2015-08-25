package service
import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/network/connection"
	"github.com/joernweissenborn/aurarath/messages"
	"errors"
)

type ServiceConnection struct {

	uuid string

	connections map[string]eventual2go.StreamController

	connected *eventual2go.Future
	disconnected *eventual2go.Future

	handshake *eventual2go.Future
	codecs []byte

}

func NewServiceConnection(uuid string) (sc *ServiceConnection){
	sc = new(ServiceConnection)
	sc.uuid = uuid
	sc.connections = map[string]eventual2go.StreamController{}
	sc.connected = eventual2go.NewFuture()
	sc.disconnected = eventual2go.NewFuture()
	sc.handshake = eventual2go.NewFuture()

	return
}

func (sc *ServiceConnection) Connected() *eventual2go.Future{
	return sc.connected
}

func (sc *ServiceConnection) Disconnected() *eventual2go.Future{
	return sc.disconnected
}

func (sc *ServiceConnection) Handshake() *eventual2go.Future{
	return sc.handshake
}

func (sc *ServiceConnection) Connect(name, address string, port int) {
	var err error
	if _, f := sc.connections[address]; !f{
		sc.connections[address], err = connection.NewOutgoing(name, address, port)
		if err != nil {
			delete(sc.connections,address)
		} else if !sc.connected.IsComplete() {
			sc.connected.Complete(sc)
		}
	}
}



func (sc *ServiceConnection) ShakeHand(codecs []byte){
	sc.codecs = codecs
	sc.handshake.Complete(sc)
}


func (sc *ServiceConnection) DoHandshake(codecs []byte,address string, port int){
	m := &messages.Hello{codecs,address,port}
	sc.Send(messages.Flatten(m))
}

func (sc *ServiceConnection) DoHandshakeReply(codecs []byte){
	m := &messages.HelloOk{codecs}
	sc.Send(messages.Flatten(m))
}


func (sc *ServiceConnection) Disconnect(addr string) {
	if conn, f := sc.connections[addr]; f {
		conn.Close()
		delete(sc.connections, addr)
	}
	if len(sc.connections) == 0 {
		sc.disconnected.Complete(sc.uuid)
	}
}

func (sc *ServiceConnection) DisconnectAll() {
	for addr, _ := range sc.connections {
		sc.Disconnect(addr)
	}
}

func (sc *ServiceConnection) Send(msg [][]byte) (err error) {
	if !sc.connected.IsComplete(){
		return errors.New("Not Connected")
	}
	for _, conn := range sc.connections {
		conn.Add(msg)
		return
	}
	return
}