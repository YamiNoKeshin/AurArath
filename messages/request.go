package messages

import (
	"bytes"
	"encoding/gob"
	uuid "github.com/nu7hatch/gouuid"
	"strings"
)

type Request struct {
	UUID     string
	Importer string
	Function string
	CallType CallType
	Codec    Codec
	params   []byte
}

func NewRequest(importer, function string, call_type CallType, params []byte) (r *Request) {
	r = new(Request)
	id, _ := uuid.NewV4()
	r.UUID = id.String()
	r.Importer = importer
	r.Function = function
	r.CallType = call_type
	r.params = params
	return
}

func (*Request) GetType() MessageType { return REQUEST }

func (r *Request) Unflatten(d []string) {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(r)
	r.params = []byte(d[1])
}

func (r *Request) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(r)
	return [][]byte{payload.Bytes(), r.params}
}

func (r *Request) Parameter() []byte {
	return r.params
}
