package messages
import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/network/connection"
)

type IncomingMessage struct {
	Sender string
	Iface string
	Msg Message
}

func ToIncomingMsg(d eventual2go.Data) eventual2go.Data {
	m := d.(connection.Message)
	return IncomingMessage{m.Payload[0],m.Iface,Unflatten(m.Payload[2:])}
}