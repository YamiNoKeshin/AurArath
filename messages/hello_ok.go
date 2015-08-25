package messages
import (
	"strings"
	"encoding/gob"
	"bytes"
)


type HelloOk struct {
	Codecs []byte
}

func(*HelloOk) GetType() MessageType {return HELLO_OK}


func(h *HelloOk) Unflatten(d []string)  {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(&h)
}

func(h *HelloOk) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(h)
	return [][]byte{payload.Bytes()}
}

