package messages
import (
	"github.com/joernweissenborn/aurarath/network/peer"
	"strings"
	"encoding/gob"
	"bytes"
)


type HelloOk struct {
	PeerDetails peer.Details
}

func(*HelloOk) GetType() MessageType {return HELLO_OK}


func(h *HelloOk) Unflatten(d []string)  {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(&h.PeerDetails)
}

func(h *HelloOk) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(h.PeerDetails)
	return [][]byte{payload.Bytes()}
}

