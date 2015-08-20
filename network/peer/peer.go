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

}

func New(uuid string) (p *Peer) {

	p = new(Peer)

	p.uuid = uuid

	p.connected = eventual2go.NewFuture()

	p.connections = map[string]eventual2go.StreamController{}

	return
}

func (p *Peer) OpenConnection(ip string, port uint16){
	var err error

	if _, f := p.connections[ip]; !f{
		p.connections[ip], err = connection.NewOutgoing(p.uuid,ip,int(port))
		if err != nil {
			delete(p.connections,ip)
		}
	}

}

func (p *Peer) isConnected() bool {
	return len(p.connections) != 0
}

func (p *Peer) Send(msg [][]byte) (err error) {
	if !p.isConnected(){
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
}

func (p *Peer) closeAllConnections() {
	for addr, conn := range p.connections {
		conn.Close()
		delete(p.connections, addr)
	}
}
