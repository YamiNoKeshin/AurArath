package daemon

import (
	"github.com/joernweissenborn/aurarath/network/connection"
	"github.com/joernweissenborn/aurarath/service"
	"github.com/joernweissenborn/eventual2go"
	"log"
	"os"
)

type Daemon struct {
	clients   map[string]*eventual2go.StreamController
	announcer map[string]*service.Announcer

	logger *log.Logger
}

func New(addr string, port int) (d *Daemon) {

	d = new(Daemon)

	d.clients = map[string]*eventual2go.StreamController{}
	d.announcer = map[string]*service.Announcer{}

	d.logger = log.New(os.Stdout, "daemon ", log.Lshortfile)
	d.logger.Println("starting up")

	incoming, _ := connection.NewIncoming(addr)
	d.logger.Println("launched incoming at", incoming.Port())

	msg := incoming.In().Where(ValidMessage)

	ct := NewClientTracker(addr, port)

	r := eventual2go.NewReactor()

	r.React("new_client", d.newClient(incoming.Port()))
	r.AddStream("new_client", ct.new.Stream)

	r.React("client_gone", d.clientGone)
	r.AddStream("client_gone", ct.gone.Stream)

	r.React("client_service", d.clientService)
	r.AddStream("client_service", msg.Where(IsMessage(EXPORT)).Transform(ToNewServiceMessage))
	r.AddStream("client_service", msg.Where(IsMessage(IMPORT)).Transform(ToNewServiceMessage))

	d.logger.Println("starting tracker")
	ct.Run()

	d.logger.Println("started")
	return
}

func (d *Daemon) newClient(port int) eventual2go.Subscriber {
	return func(data eventual2go.Data) {
		np := data.(Newclient)
		d.logger.Println("new client", np)
		conn, err := connection.NewOutgoing("AURARATH_DAEMON", "127.0.0.1", np.Port)
		if err != nil {
			d.logger.Println("error opening connection:", err)
			return
		}
		conn.Add(NewHello(port))
		d.clients[np.UUID] = conn

	}
}

func (d *Daemon) clientGone(data eventual2go.Data) {
	gp := data.(string)
	d.clients[gp].Close()
}

func (d *Daemon) clientService(data eventual2go.Data) {
	m := data.(NewService)
	d.logger.Println("new service", m)
	a := service.NewAnnouncer(m.UUID, m.Addresses, m.ServiceType, m.Descriptor)
	client := d.clients[m.UUID]
	a.ServiceArrived().Listen(serviceArrived(client))
	a.ServiceGone().Listen(serviceGone(client))
	a.Run()
	d.announcer[m.UUID] = a
}

func serviceArrived(client *eventual2go.StreamController) eventual2go.Subscriber {
	return func(d eventual2go.Data) {
		sa := d.(service.ServiceArrived)
		client.Add(NewServiceArrived(sa))
	}
}
func serviceGone(client *eventual2go.StreamController) eventual2go.Subscriber {
	return func(d eventual2go.Data) {
		uuid := d.(string)
		client.Add(NewServiceGone(uuid))
	}
}
