package daemon
import (
	"time"
	"github.com/joernweissenborn/eventual2go"
	"net"
	"encoding/binary"
)

type PeerTracker struct {
	peers map[string]time.Time
	new eventual2go.StreamController
	gone eventual2go.StreamController
}

func NewPeerTracker() (p *PeerTracker) {
	p = new(PeerTracker)
	p.peers = map[string]time.Time
	p.new = eventual2go.NewStreamController()
	p.gone = eventual2go.NewStreamController()
	return
}

func (p *PeerTracker) Run() {
	r := eventual2go.NewReactor()
	r.React("ping",p.checkInPeer)
	r.AddStream("ping",listenUdp())
	r.React("check",p.checkPeers)
	go func() {
		for _ := range time.Tick(1*time.Second) {
			r.Fire("check",nil)
		}
		}()
	}

}

func (p *PeerTracker) New() eventual2go.Stream {
	return p.new.Stream
}

func (p *PeerTracker) Gone() eventual2go.Stream {
	return p.gone.Stream
}

func (p *PeerTracker) checkInPeer(d eventual2go.Data) {
	sig := d.([]byte)
	id := string(sig[:32])
	if _, f := p.peers;!f {
		port := binary.LittleEndian.Uint16(sig[32:])
		p.new.Add(NewPeer{id,port})
	}
	p.peers[id] = time.Now()
}

func (p *PeerTracker) checkPeers(d eventual2go.Data) {
	for uuid, t := range p.peers {
		if time.Since(t)>3*time.Second {
			p.gone.Add(uuid)
		}
	}
}


func listenUdp() eventual2go.Stream{
	conn, err := net.ListenUDP("ipv4",&net.UDPAddr{IP: net.ParseIP("127.0.0.1"),Port:5558})
	if err != nil {
		panic(err)
	}
	s := eventual2go.NewStreamController()
	go func(){
		for {
			data := make([]byte, 128)
			read, addr, _ := conn.ReadFromUDP(data)
			s.Add(data[:read])
		}
	}()
	return s.Stream
}

type NewPeer struct {
	UUID string
	Port int
}