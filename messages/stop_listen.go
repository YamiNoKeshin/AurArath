package messages

import (
	"bytes"
	"encoding/gob"
	"strings"
)

type StopListen struct {
	Function string
}

func (*StopListen) GetType() MessageType { return STOP_LISTEN }

func (l *StopListen) Unflatten(d []string) {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(l)
}

func (l *StopListen) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(l)
	return [][]byte{payload.Bytes()}
}
