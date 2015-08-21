package messages
import (
	"encoding/gob"
	"github.com/joernweissenborn/eventual2go"
)



//go:generate stringer -type=MessageType

type MessageType int

const (
	HELLO MessageType = iota
	HELLO_OK
)

func Get(messagetype MessageType) (msg Message){

	switch messagetype {
	case HELLO:
		msg =new(Hello)
	case HELLO_OK:
		msg =new(HelloOk)
	}
	return
}

func init(){
	gob.Register(Hello{})
	gob.Register(HelloOk{})
}

func Is(t MessageType) eventual2go.Filter {
	return func(d eventual2go.Data) bool {return d.(IncomingMessage).Msg.GetType() == t}
}
