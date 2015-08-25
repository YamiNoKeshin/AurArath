package messages
import (
	"encoding/gob"
	"github.com/joernweissenborn/eventual2go"
	"log"
)



//go:generate stringer -type=MessageType

type MessageType int

const (
	HELLO MessageType = iota
	HELLO_OK
	REQUEST
	RESULT
	LISTEN
	STOP_LISTEN
)

func Get(messagetype MessageType) (msg Message){

	switch messagetype {
	case HELLO:
		msg =new(Hello)
	case HELLO_OK:
		msg =new(HelloOk)
	case REQUEST:
		msg=new(Request)
	case RESULT:
		msg=new(Result)
	case LISTEN:
		msg=new(Listen)
	case STOP_LISTEN:
		msg=new(StopListen)
	}
	return
}

func init(){
	gob.Register(Hello{})
	gob.Register(HelloOk{})
	gob.Register(Request{})
	gob.Register(Result{})
}

func Is(t MessageType) eventual2go.Filter {
	return func(d eventual2go.Data) bool {
		log.Println("IS",d)
		return d.(IncomingMessage).Msg.GetType() == t}
}
