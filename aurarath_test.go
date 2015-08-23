package aurarath
import (
	"testing"
	"time"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/messages"
	"bytes"
)

func TestCallOne2One(t *testing.T) {
	i := NewImport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer i.Remove()
	e := NewExport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer e.Remove()
	c := e.Requests().AsChan()
	i.Run()

	e.Run()

	params := []byte{4,5,63,4}
	f := i.Call("SeyHello",params)

	select {
	case <-time.After(2*time.Second):
		t.Fatal("Didnt Got Request")
	case d := <-c:
		r := d.(*messages.Request)
		if r.Importer != i.UUID() {
			t.Error("Wrong Import UUID",r.Importer,i.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params",r.Parameter(),params)
		}
		e.Reply(r,params)
	}


	select {
	case <-time.After(2*time.Second):
		t.Fatal("Didnt Got Request")
	case d := <-f.AsChan():
		r := d.(*messages.Result)
		if r.Exporter != e.UUID() {
			t.Error("Wrong Export UUID",r.Exporter,e.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params",r.Parameter(),params)
		}
	}


}

var TEST_APP_DESCRIPTOR *AppDescriptor = &AppDescriptor{
	[]Function{
		Function{
			"SayHello",
			[]Parameter{
				Parameter{
					"Greeting",
					"string"}},
			[]Parameter{
				Parameter{
					"Answer",
					"string"}},
		},
	},
	map[string]string{"TAG_1":""},
}