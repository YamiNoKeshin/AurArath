package messages
import (
	uuid "github.com/nu7hatch/gouuid"
	"strings"
	"encoding/gob"
	"bytes"
)


type Result struct {
	Request *Request
	Exporter string
	result []byte
}

func NewResult(Export, Request, Result []byte) (r *Result) {
	r = new(Result)
	r.Exporter= Export
	r.Request = Request
	r.result = Result
	return
}

func(*Result) GetType() MessageType {return Result}


func(r *Result) Unflatten(d []string)  {
	dec := gob.NewDecoder(strings.NewReader(d[0]))
	dec.Decode(r)
	r.result = []byte(d[1])
}

func(r *Result) Flatten() [][]byte {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)
	enc.Encode(r)
	return [][]byte{payload.Bytes(),r.result}
}

