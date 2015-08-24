package daemon
import (
	"github.com/joernweissenborn/aurarath/network/beacon"
	"net"
	"github.com/joernweissenborn/eventual2go"
	"time"
	"github.com/joernweissenborn/aurarath/network/connection"
)


type Daemon struct {
	peers map[string] eventual2go.StreamController
}

func New() (d *Daemon) {

	d = new(Daemon)
	incoming,_ := connection.NewIncoming("127.0.0.1")
	msg := incoming.In().Where(ValidMessage)
	r := eventual2go.NewReactor()
	pt := NewPeerTracker()
	r.React("new_peer",d.newPeer(incoming.Port()))
	r.AddStream("new_peer",pt.new)
	r.React("gone_peer",d.peerGone)
	r.AddStream("new_peer",pt.gone)
	r.React("peer_export",d.peerExport)
	r.AddStream("peer_export",msg.Where(IsMessage(EXPORT)).Transform(ToNewServiceMessage))
	r.React("peer_import",d.peerImport)
	r.AddStream("peer_import",msg.Where(IsMessage(IMPORT)).Transform(ToNewServiceMessage))
}

func (d *Daemon) newPeer(port int)eventual2go.Subscription {
	return func(d eventual2go.Data) {
		np := d.(NewPeer)
		conn,err := connection.NewOutgoing("AURARATH_DAEMON","127.0.0.1",np.Port)
		if err != nil {
			return
		}
		conn.Add(NewHello(port))
		d.peers[np.UUID] = conn

	}
}

func (d *Daemon) peerGone(d eventual2go.Data) {
	gp := d.(string)
	d.peers[gp].Close()
}


func (d *Daemon) peerExport(d eventual2go.Data) {
}

func (d *Daemon) peerImport(d eventual2go.Data) {
}


