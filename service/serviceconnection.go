package service

import (
	"errors"
	"github.com/joernweissenborn/aurarath/messages"
	"github.com/joernweissenborn/aurarath/network/connection"
	"github.com/joernweissenborn/eventual2go"
)

type ServiceConnection struct {
	uuid string

	connections map[string]*eventual2go.StreamController

	connected    *eventual2go.Completer
	disconnected *eventual2go.Completer

	handshake *eventual2go.Completer
	codecs    []byte

	handshaked bool
}

func NewServiceConnection(uuid string) (sc *ServiceConnection) {
	sc = new(ServiceConnection)
	sc.uuid = uuid
	sc.connections = map[string]*eventual2go.StreamController{}
	sc.connected = eventual2go.NewCompleter()
	sc.disconnected = eventual2go.NewCompleter()
	sc.handshake = eventual2go.NewCompleter()

	return
}

func (sc *ServiceConnection) Uuid() string {
	return sc.uuid
}
func (sc *ServiceConnection) Connected() *eventual2go.Future {
	return sc.connected.Future()
}

func (sc *ServiceConnection) Disconnected() *eventual2go.Future {
	return sc.disconnected.Future()
}

func (sc *ServiceConnection) Handshake() *eventual2go.Future {
	return sc.handshake.Future()
}

func (sc *ServiceConnection) Connect(name, address string, port int) {
	var err error
	if _, f := sc.connections[address]; !f {
		sc.connections[address], err = connection.NewOutgoing(name, address, port)
		if err != nil {
			delete(sc.connections, address)
		} else if !sc.connected.Completed() {
			sc.connected.Complete(sc)
		}
	}
}

func (sc *ServiceConnection) ShakeHand(codecs []byte) {
	sc.codecs = codecs
	sc.handshake.Complete(sc)
}

func (sc *ServiceConnection) DoHandshake(codecs []byte, address string, port int) {
	m := &messages.Hello{codecs, address, port}
	sc.Send(messages.Flatten(m))
}

func (sc *ServiceConnection) DoHandshakeReply(codecs []byte) {
	if sc.handshaked {
		return
	}
	m := &messages.HelloOk{codecs}
	sc.Send(messages.Flatten(m))
	sc.handshaked = true
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
	if !sc.connected.Completed() {
		return errors.New("Not Connected")
	}
	for _, conn := range sc.connections {
		conn.Add(msg)
		return
	}
	return
}
