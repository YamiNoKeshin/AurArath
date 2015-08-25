package messages
import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/network"
	"strconv"
	"github.com/joernweissenborn/aurarath/network/connection"
	"log"
)



type Message interface {
	GetType() MessageType
	Flatten() [][]byte
	Unflatten([]string)
}

func Flatten(m Message) [][]byte {
	t := strconv.FormatInt(int64(m.GetType()),10)
	payload := [][]byte{[]byte{byte(network.PROTOCOL_SIGNATURE)},[]byte(t)}
	for _,p := range m.Flatten() {
		payload = append(payload,p)
	}
	return payload
}

func Unflatten(m []string) (msg Message){
	mtype, _ := strconv.ParseInt(m[0],10,8)
	msg = Get(MessageType(mtype))
	msg.Unflatten(m[1:])
	return
}

func Valid(d eventual2go.Data) bool {
	log.Println("VALID",d)
	m := d.(connection.Message).Payload
	if len(m)<3 {
		return false
	}
	p := []byte(m[1])[0]


	if p != network.PROTOCOL_SIGNATURE{
		return false
	}
	return true
}