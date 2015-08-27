package messages

import (
	"github.com/joernweissenborn/aurarath/network/connection"
	"github.com/joernweissenborn/eventual2go"
)

type IncomingMessage struct {
	Sender string
	Iface  string
	Msg    Message
}

func ToIncomingMsg(d eventual2go.Data) eventual2go.Data {
	m := d.(connection.Message)
	return IncomingMessage{m.Payload[0], m.Iface, Unflatten(m.Payload[2:])}
}

func ToMsg(d eventual2go.Data) eventual2go.Data {
	return d.(IncomingMessage).Msg
}
