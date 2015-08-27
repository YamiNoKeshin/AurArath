package messages

import (
	"bytes"
	"encoding/gob"
	"strings"
)

type Hello struct {
	Codecs  []byte
	Address string
	Port    int
}

func (*Hello) GetType() MessageType { return HELLO }

func (h *Hello) Unflatten(d []string) {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(&h)
}

func (h *Hello) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(h)
	return [][]byte{payload.Bytes()}
}
