package messages

import (
	"bytes"
	"encoding/gob"
	"strings"
)

type Listen struct {
	Function string
}

func (*Listen) GetType() MessageType { return LISTEN }

func (l *Listen) Unflatten(d []string) {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(l)
}

func (l *Listen) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(l)
	return [][]byte{payload.Bytes()}
}
