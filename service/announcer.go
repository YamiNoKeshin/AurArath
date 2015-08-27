package service

import (
	"github.com/joernweissenborn/aurarath/network/node"
	"strings"
	"strconv"
	"github.com/joernweissenborn/eventual2go"
	"encoding/binary"
	"log"
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"github.com/joernweissenborn/aurarath/config"
	"fmt"
)

type Announcer struct{
	node *node.Node
	r *eventual2go.Reactor
	servicetype string

	clientPorts map[string]int

	new *eventual2go.StreamController
	logger *log.Logger

	announced *eventual2go.Completer
}

func NewAnnouncer(uuid string, addresses []string, servicetype string, desc *appdescriptor.AppDescriptor) (a *Announcer){

	cfg := config.DefaultLocalhost()

	a = new(Announcer)
	a.announced = eventual2go.NewCompleter()
	a.logger = log.New(cfg.Logger(),fmt.Sprintf("announcer %s ",uuid),log.Lshortfile)

	a.new = eventual2go.NewStreamController()
	a.servicetype = servicetype
	addrs := []string{}
	a.clientPorts = map[string]int{}

	for _, addr := range addresses {
		as := strings.Split(addr,":")
		addrs = append(addrs,as[0])
		p, _ := strconv.ParseInt(as[1],0,0)
		a.clientPorts[as[0]] = int(p)
		a.logger.Println("adding address",as[0],int(p))
	}

	cfg.NetworkInterfaces = addrs
	a.node = node.New(uuid,cfg,desc.AsTagSet())

	a.r = eventual2go.NewReactor()
	a.r.React("first_join",a.announce)
	a.r.AddFuture("first_join",a.node.Join().First())
	a.r.React("service_found",a.replyToServiceQuery)
	a.r.AddStream("service_found",a.node.Queries().WhereNot(isService(a.servicetype)))
	a.logger.Println("setup finished")
	return
}

func (a *Announcer) Announced() *eventual2go.Future{
	return a.announced.Future()
}

func (a *Announcer) ServiceArrived() *eventual2go.Stream {
	return a.new.Stream
}

func (a *Announcer) ServiceGone() *eventual2go.Stream {
	return a.node.Leave().Transform(toServiceGone)
}

func (a *Announcer) Run() {
	a.logger.Println("starting")
	a.node.Run()
}

func (a *Announcer) Shutdown() {
	a.node.Shutdown()
	a.r.Shutdown()
}

func (a *Announcer) announce(eventual2go.Data) {
	a.logger.Println("announcing")
	results := eventual2go.NewStreamController()
	c := results.AsChan()
	a.node.Query(a.servicetype,nil,results)
	go a.collectAnnounceResponses(c)
	return
}

func (a *Announcer) collectAnnounceResponses(c chan eventual2go.Data) {
	for d := range c {
		r := d.(node.QueryResponseEvent)
		buf := strings.Split(r.Response.From, "@")
		if len(buf) != 2 {
			return
		}

		uuid := buf[0]
		ip := buf[1]
		port := binary.LittleEndian.Uint16(r.Response.Payload)
		a.logger.Printf("got reply from %s at %s:%d",uuid,ip,port)
		a.new.Add(ServiceArrived{uuid,r.Interface,ip,int(port)})
	}
	a.logger.Printf("finished announce")
	a.announced.Complete(nil)
}

func (a *Announcer) replyToServiceQuery(d eventual2go.Data){
	q := d.(node.QueryEvent)
	a.logger.Println("found service on",q.Address)
	if port, f := a.clientPorts[q.Address];f {
		repl := make([]byte,2)
		binary.LittleEndian.PutUint16(repl,uint16(port))
		q.Query.Respond(repl)
	}
}


type ServiceArrived struct {
	UUID string
	Interface string
	Address string
	Port int
}

type ServiceGone struct {
	UUID string
	Address string
}

func toServiceGone(d eventual2go.Data) eventual2go.Data {
	buf := strings.Split(d.(node.LeaveEvent).Name, "@")
	return ServiceGone{buf[0],buf[1]}
}