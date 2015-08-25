package service

import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/network/connection"
	"log"
	"github.com/joernweissenborn/aurarath/messages"
	"fmt"
	"github.com/joernweissenborn/aurarath/appdescriptor"
	uid "github.com/nu7hatch/gouuid"
)


type Service struct {

	uuid string

	r *eventual2go.Reactor

	appDescriptor *appdescriptor.AppDescriptor

	announcer *Announcer

	remove *eventual2go.Future

	incoming map[string]*connection.Incoming

	in eventual2go.StreamController

	connectedServices map[string]*ServiceConnection

	servicetype string

	codecs []byte

	logger *log.Logger

	connected *eventual2go.Future
	disconnected *eventual2go.Future

	newpeers eventual2go.StreamController
	gonepeers eventual2go.StreamController
}

func NewService(a *appdescriptor.AppDescriptor, servicetype string, cfg *config.Config, codecs []byte) (s *Service){
	s = new(Service)
	id,_:= uid.NewV4()
	s.uuid = id.String()
	s.logger = log.New(cfg.Logger(),fmt.Sprintf("service %s  ",id),log.Lshortfile)
	s.appDescriptor = a
	s.servicetype = servicetype
	s.codecs = codecs
	s.newpeers = eventual2go.NewStreamController()
	s.gonepeers = eventual2go.NewStreamController()

	s.incoming = map[string]*connection.Incoming{}
	s.in = eventual2go.NewStreamController()

	s.connectedServices = map[string]*ServiceConnection{}
	s.connected = eventual2go.NewFuture()
	s.disconnected = eventual2go.NewFuture()
	s.remove = eventual2go.NewFuture()

	s.r = eventual2go.NewReactor()
	s.r.React("service_arrived",s.serviceArrived)
	s.r.React("service_gone",s.serviceGone)
	s.r.React("service_shake_hand",s.serviceHandshake)
	s.r.AddStream("service_shake_hand",s.in.Where(messages.Is(messages.HELLO)))
	s.r.React("service_shake_hand_reply",s.serviceHandShakeReply)
	s.r.AddStream("service_shake_hand_reply",s.in.Where(messages.Is(messages.HELLO_OK)))
	s.createIncoming(cfg)
	s.createAnnouncer()
	return
}

func (s *Service) UUID() string{
	return s.uuid
}

func (s *Service) Connected() *eventual2go.Future{
	return s.connected
}

func (s *Service) Disconnected() *eventual2go.Future{
	return s.disconnected
}

func (s *Service) Run() {
	s.announcer.Run()
}

func (s *Service) createIncoming(cfg *config.Config) {
	for _, addr := range cfg.NetworkInterfaces {
		s.logger.Println("Opening Incoming Socket on", addr)
		incoming, err := connection.NewIncoming(addr)
		if err == nil {
			s.in.Join(incoming.In().Where(messages.Valid).Transform(messages.ToIncomingMsg))
			s.logger.Println("port is", incoming.Port())
			s.incoming[addr] = incoming
		} else {
			s.logger.Println("Error opening socket",err)
		}
	}
}

func (s *Service) createAnnouncer() {
	addrs := []string{}

	for addr,i := range s.incoming {
		s.logger.Println(addr,i)
		addrs = append(addrs,fmt.Sprintf("%s:%d",addr,i.Port()))
	}

	s.announcer = NewAnnouncer(s.uuid,addrs,s.servicetype,s.appDescriptor)
	s.r.AddStream("service_arrived",s.announcer.ServiceArrived())
	s.r.AddStream("service_gone",s.announcer.ServiceGone())
}

func (s *Service) serviceArrived(d eventual2go.Data) {
	sa := d.(ServiceArrived)
	s.logger.Println("Service arrived at",sa.Address,sa.Port)
	if !s.serviceConnectionExists(sa.UUID) {
		s.logger.Println("Service does not exist, creating",sa.UUID)
		sc := s.createServiceConnection(sa.UUID)
		sc.Connect(s.UUID(),sa.Address,sa.Port)
		sc.Connected().Then(s.doHandShake(sa.Interface))
	}
}

func (s *Service) serviceGone(d eventual2go.Data) {
	r := d.(ServiceGone)

	if sc, f:= s.connectedServices[r.UUID];f{
		sc.Disconnect(r.Address)
	}
}

func (s *Service) createServiceConnection(uuid string) (sc *ServiceConnection) {
	sc = NewServiceConnection(uuid)
	sc.Disconnected().Then(s.removeServiceConnection)
	s.connectedServices[uuid] = sc
	return
}
func (s *Service) removeServiceConnection(d eventual2go.Data) (eventual2go.Data) {
	uuid := d.(string)
	s.logger.Println("Removing service connection",uuid)
	delete(s.connectedServices,uuid)
	if len(s.connectedServices) == 0 && !s.disconnected.IsComplete() {
		s.logger.Println("Disconnected")
		s.disconnected.Complete(nil)
	}
	return nil
}

func (s *Service) serviceHandshake(d eventual2go.Data) {
	m := d.(messages.IncomingMessage)
	h := m.Msg.(*messages.Hello)

	s.logger.Println("Got handshake:",m.Sender,h.Address,h.Port)

	sc := s.connectedServices[m.Sender]
	if sc == nil {
		s.logger.Println("Service does not exist, creating",m.Sender)
		sc = s.createServiceConnection(m.Sender)
	}
	sc.Connected().Then(s.doHandShakeReply())

	sc.Connect(s.UUID(), h.Address,h.Port)
	s.logger.Println("DONE")

}


func (s *Service) doHandShake(iface string) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		sc := d.(*ServiceConnection)
		s.logger.Println("doing handshake with",sc.uuid)
		port := s.incoming[iface].Port()
		sc.DoHandshake(s.codecs,iface,port)
		return nil
	}
}

func (s *Service) doHandShakeReply() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data{
		sc := d.(*ServiceConnection)
		s.logger.Println("replying handshake to",sc.uuid)
		sc.DoHandshakeReply(s.codecs)
		return nil
	}
}

func (s *Service) serviceHandShakeReply(d eventual2go.Data) {
	m := d.(messages.IncomingMessage)
	h := m.Msg.(*messages.HelloOk)
	s.logger.Println("Got handshake reply from",m.Sender)
	sc := s.connectedServices[m.Sender]

	sc.ShakeHand(h.Codecs)
	if !s.connected.IsComplete() {
		s.logger.Println("Connected")
		s.connected.Complete(m.Sender)
	}
}

func (s *Service) serviceConnectionExists(uuid string) (e bool){
	_,e = s.connectedServices[uuid]
	return
}
func (s *Service) GetConnectedService(uuid string) (sc *ServiceConnection){
	return s.connectedServices[uuid]
}

func (s *Service) GetConnectedServices() (scs []*ServiceConnection){
	scs = []*ServiceConnection{}
	for _,sc := range s.connectedServices{
		if sc.Connected().IsComplete() && sc.Handshake().IsComplete() {
			scs = append(scs,sc)
		}
	}
	return
}


func (s *Service) Remove() {
	s.logger.Println("Stopping Service",s.UUID())
	s.announcer.Shutdown()
	for _, i := range s.incoming {
		i.Close()
	}
	for _, sc := range s.connectedServices{
		sc.DisconnectAll()
	}

	s.remove.Complete(nil)
	s.logger.Println("Service Stopped",s.UUID())
}


