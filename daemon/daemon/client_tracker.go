package daemon
import (
	"time"
	"github.com/joernweissenborn/eventual2go"
	"net"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

type ClientTracker struct {

	address string
	port int

	clients map[string]time.Time

	new *eventual2go.StreamController
	gone *eventual2go.StreamController

	logger *log.Logger
}

func NewClientTracker(address string, port int) (p *ClientTracker) {
	p = new(ClientTracker)
	p.address = address
	p.port = port
	p.clients = map[string]time.Time{}
	p.new = eventual2go.NewStreamController()
	p.gone = eventual2go.NewStreamController()
	p.logger = log.New(os.Stdout,"clienttracker ",log.Lshortfile)
	return
}

func (p *ClientTracker) Run() {
	p.logger.Println("starting")
	r := eventual2go.NewReactor()
	r.React("ping",p.checkInclient)
	r.AddStream("ping",p.listenUdp())
	r.React("check",p.checkclients)
	go func(reactor *eventual2go.Reactor) {
		for range time.Tick(1*time.Second) {
			reactor.Fire("check",nil)
		}
	}(r)
}

func (p *ClientTracker) New() *eventual2go.Stream {
	return p.new.Stream
}

func (p *ClientTracker) Gone() *eventual2go.Stream {
	return p.gone.Stream
}

func (p *ClientTracker) checkInclient(d eventual2go.Data) {
	sig := strings.Split(string(d.([]byte)),":")
	id := sig[0]
	if _, f := p.clients[id];!f {
		port,_ := strconv.ParseInt(sig[1],0,0)
		p.new.Add(Newclient{id,int(port)})
	}
	p.clients[id] = time.Now()
}

func (p *ClientTracker) checkclients(d eventual2go.Data) {
	for uuid, t := range p.clients {
		if time.Since(t)>3*time.Second {
			p.gone.Add(uuid)
		}
	}
}


func (p *ClientTracker) listenUdp() *eventual2go.Stream{
	addr := &net.UDPAddr{IP: net.ParseIP(p.address),Port:p.port}
	p.logger.Println("Starting to listen on",addr)
	conn, err := net.ListenUDP("udp4",addr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	s := eventual2go.NewStreamController()
	go func(stream *eventual2go.StreamController){
		for {
			data := make([]byte, 128)
			read, _, _ := conn.ReadFromUDP(data)
			stream.Add(data[:read])
		}
	}(s)

	return s.Stream
}

type Newclient struct {
	UUID string
	Port int
}