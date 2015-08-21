package peer
import (
	"github.com/joernweissenborn/eventual2go"
	"sync"
	"github.com/joernweissenborn/aurarath/network/connection"
	"errors"
)

type Details struct {
	Codecs []byte
}

type Peer struct {

	m *sync.Mutex

	uuid string

	details Details

	connections map[string]eventual2go.StreamController

	connected *eventual2go.Future
	disconnected *eventual2go.Future

	greeted *eventual2go.Future

}

func New(uuid string) (p *Peer) {

	p = new(Peer)

	p.uuid = uuid

	p.connected = eventual2go.NewFuture()
	p.disconnected = eventual2go.NewFuture()
	p.greeted = eventual2go.NewFuture()

	p.connections = map[string]eventual2go.StreamController{}

	return
}

func (p *Peer) Uuid() string {
	return p.uuid
}

func (p *Peer) OpenConnection(ip string, port uint16, id string){
	var err error
	if _, f := p.connections[ip]; !f{
		p.connections[ip], err = connection.NewOutgoing(id, ip, int(port))
		if err != nil {
			delete(p.connections,ip)
		} else if !p.connected.IsComplete() {
			p.connected.Complete(p)
		}
	}

}

func (p *Peer) Connected() *eventual2go.Future {
	return p.connected
}
func (p *Peer) Disconnected() *eventual2go.Future {
	return p.disconnected
}

func (p *Peer) Send(msg [][]byte) (err error) {
	if !p.connected.IsComplete(){
		return errors.New("Not Connected")
	}

	for _, conn := range p.connections {
		conn.Add(msg)
		return
	}
	return
}

func (p *Peer) Details()(d Details) {
	return p.details
}

func (p *Peer) SetDetails(d Details) {
	p.details = d
	p.greeted.Complete(nil)
}

func (p *Peer) Greeted()(*eventual2go.Future) {
	return p.greeted
}

func (p *Peer) CloseConnection(addr string) {
	if conn, f := p.connections[addr]; f {
		conn.Close()
		delete(p.connections, addr)
	}
	if len(p.connections) == 0 {
		p.disconnected.Complete(p.uuid)
	}
}

func (p *Peer) CloseAllConnections() {
	for addr, _ := range p.connections {
		p.CloseConnection(addr)
	}
}
