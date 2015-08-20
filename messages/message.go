package messages
import (
	"encoding/gob"
	"encoding/binary"
	"bytes"
)



type Message interface {
	GetType() uint16
}

func Flatten(m Message) [][]byte {
	t := make([]byte,2)
	binary.LittleEndian.PutUint16(t,m.GetType())
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(m)
	return [][]byte{t,payload.Bytes()}
}