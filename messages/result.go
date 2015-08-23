package messages
import (
	"strings"
	"encoding/gob"
	"bytes"
)


type Result struct {
	Request *Request
	Exporter string
	params []byte
}

func NewResult(export string, request *Request, parameter []byte) (r *Result) {
	r = new(Result)
	r.Exporter= export
	r.Request = request
	r.params = parameter
	return
}

func(*Result) GetType() MessageType {return RESULT}


func(r *Result) Unflatten(d []string)  {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(r)
	r.params = []byte(d[1])
}

func(r *Result) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(r)
	return [][]byte{payload.Bytes(),r.params}
}


func(r *Result) Parameter() []byte {
	return r.params
}

