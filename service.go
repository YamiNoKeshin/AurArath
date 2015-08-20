package connection

import (
	"github.com/joernweissenborn/aurarath/network/node"
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/network/connection"
	"encoding/binary"
	"strings"
	"log"
	"github.com/joernweissenborn/aurarath/network/peer"
	"github.com/joernweissenborn/aurarath/messages"
)


type Service struct {
	appDescriptor *AppDescriptor

	node   *node.Node

	remove *eventual2go.Future

	incoming map[string]connection.Incoming
	peers map[string]*peer.Peer
	
	servicetype string

	codecs []byte

	logger *log.Logger
}

func NewService(a *AppDescriptor, servicetype string, cfg *config.Config, codecs []byte) (s *Service){
	s = new(Service)
	s.logger = log.New(cfg.Logger(),"Service ",log.Lshortfile)
	s.appDescriptor = a
	s.servicetype = servicetype
	s.codecs = codecs
	s.incoming = map[string]connection.Incoming{}
	s.peers = map[string]*peer.Peer{}

	s.node = node.New(cfg,a.AsTagSet())
	s.node.Queries().WhereNot(isService(servicetype)).Listen(s.replyToService)
	s.node.Join().First().Then(s.announce)
	s.createIncoming(cfg)
	return
}

func (s *Service) Run() {
	s.node.Run()
}
func (s *Service) createIncoming(cfg *config.Config) {
	for _, addr := range cfg.NetworkInterfaces {
		s.logger.Println("Opening Incoming Socket on", addr)
		incoming, err := connection.NewIncoming(addr)
		if err == nil {
			s.incoming[addr] = incoming
		} else {
			s.logger.Println("Error opening socket",err)
		}
	}
}

func (s *Service) announce(eventual2go.Data) eventual2go.Data{
	results := eventual2go.NewStreamController()
	results.Listen(s.foundPeer)
	s.node.Query(s.servicetype,nil,results)
	return nil
}

func (s *Service) foundPeer(d eventual2go.Data) {
	r := d.(node.QueryResponseEvent)
	buf := strings.Split(r.Response.From, "@")
	if len(buf) != 2 {
		return
	}
	uuid := buf[0]
	ip := buf[1]
	port := binary.LittleEndian.Uint16(r.Response.Payload)

	p := s.peers[uuid]
	if p == nil {
		p = peer.New(uuid)
		s.peers[uuid] = p
	}

	p.OpenConnection(ip, port)
	p.Send(messages.Flatten(messages.Hello{s.getPeerDetails()}))
}

func (s *Service) getPeerDetails() peer.Details{
	return peer.Details{s.codecs}
}

func (s *Service) replyToService(d eventual2go.Data){
	q := d.(node.QueryEvent)
	s.logger.Println("Found Service on",q.Address)
	s.logger.Println(s.incoming)
	if conn, f := s.incoming[q.Address];f {
		repl := make([]byte,2)
		binary.LittleEndian.PutUint16(repl,conn.Port())
		q.Query.Respond(repl)

	}
}



