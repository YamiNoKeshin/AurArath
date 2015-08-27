package daemon

import (
	"fmt"
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"github.com/joernweissenborn/aurarath/network/connection"
	"github.com/joernweissenborn/eventual2go"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type Client struct {
	arrived *eventual2go.Stream
	gone    *eventual2go.Stream

	stop *eventual2go.Completer

	uuid        string
	ad          *appdescriptor.AppDescriptor
	servicetype string
	addresses   []string

	localAddr     string
	localPort     int
	daemonAddr    string
	daemonPort    int
	daemonPortUdp int

	daemonConn *eventual2go.StreamController

	logger *log.Logger
}

func NewClient(uuid, localAddr string, daemonAddr string, daemonPortUdp int, servicetype string, ad *appdescriptor.AppDescriptor, addresses []string) (c *Client) {

	c = new(Client)
	c.logger = log.New(os.Stdout, "client ", log.Lshortfile)
	c.logger.Println("starting")
	c.uuid = uuid
	c.ad = ad
	c.addresses = addresses
	c.servicetype = servicetype
	c.localAddr = localAddr
	c.daemonAddr = daemonAddr
	c.daemonPortUdp = daemonPortUdp

	c.stop = eventual2go.NewCompleter()

	incoming, err := connection.NewIncoming(localAddr)

	if err != nil {
		panic(err)
	}
	c.localPort = incoming.Port()
	c.logger.Println("incoming port is", c.localPort)
	c.stop.Future().Then(stopIncoming(incoming))

	msg_in := incoming.In().Where(ValidMessage)
	c.arrived = msg_in.Where(IsMessage(SERVICE_ARRIVE)).Transform(ToServiceArrivedMessage)
	c.gone = msg_in.Where(IsMessage(SERVICE_GONE)).Transform(ToServiceGone)
	msg_in.Where(IsMessage(HELLO)).Listen(c.onHello)
	return
}

func (c *Client) Run() {
	c.logger.Println("starting")
	go c.pingUdp()
}

func (c *Client) Arrived() *eventual2go.Stream {
	return c.arrived
}

func (c *Client) Gone() *eventual2go.Stream {
	return c.arrived
}

func stopIncoming(i *connection.Incoming) eventual2go.CompletionHandler {
	return func(eventual2go.Data) eventual2go.Data {
		i.Close()
		return nil
	}
}

func (c *Client) onHello(d eventual2go.Data) {
	m := d.(connection.Message)
	strport := string(m.Payload[3])
	port, _ := strconv.ParseInt(strport, 0, 0)
	c.logger.Println("got hello from daemon on port", port)
	c.daemonPort = int(port)
	c.daemonConn, _ = connection.NewOutgoing(c.uuid, c.daemonAddr, c.daemonPort)
	h := NewService{c.uuid, c.ad, c.addresses, c.servicetype}
	c.daemonConn.Add(h.flatten())
}

func (c *Client) pingUdp() {

	var pingtime = 1 * time.Second

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:0", c.localAddr))
	if err != nil {
		panic(err)
	}
	daemonAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.daemonAddr, c.daemonPortUdp))

	if err != nil {
		panic(err)
	}
	c.logger.Println("starting to pindg address", daemonAddr)
	con, err := net.DialUDP("udp", localAddr, daemonAddr)
	if err != nil {
		panic(err)
	}
	t := time.NewTimer(pingtime)
	s := c.stop.Future().AsChan()

	for {
		select {
		case <-s:
			return

		case <-t.C:
			con.Write([]byte(fmt.Sprintf("%s:%d", c.uuid, c.localPort)))
			t.Reset(pingtime)
		}
	}

}

func (c *Client) Stop() {
	c.stop.Complete(nil)
}
