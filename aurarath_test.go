package aurarath_test
import (
	"testing"
	"time"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/messages"
	"bytes"
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath"
)

func TestCallOne2One(t *testing.T) {
	i := aurarath.NewImport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer i.Remove()
	e := aurarath.NewExport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer e.Remove()
	c := e.Requests().AsChan()
	i.Run()

	e.Run()
	<-i.Connected().AsChan()
	<-e.Connected().AsChan()
	params := []byte{4,5,63,4}
	f := i.Call("SeyHello",params)

	select {
	case <-time.After(5*time.Second):
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

func TestCallMany2One(t *testing.T) {
	i1 := aurarath.NewImport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer i1.Remove()
	i2 := aurarath.NewImport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer i2.Remove()
	e := aurarath.NewExport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer e.Remove()
	c := e.Requests().AsChan()

	c1 := i1.Results().AsChan()
	c2 := i2.Results().AsChan()

	e.Run()
	i2.Run()
	i1.Run()

	time.Sleep(10*time.Second)
	i1.Listen("SayHello")
	i2.Listen("SayHello")
	time.Sleep(1*time.Second)

	params := []byte{4,5,63,4}
	i1.Trigger("SayHello",params)

	select {
	case <-time.After(5*time.Second):
		t.Fatal("Didnt Got Request")
	case d := <-c:
		r := d.(*messages.Request)
		if r.Importer != i1.UUID() {
			t.Error("Wrong Import UUID",r.Importer,i1.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params",r.Parameter(),params)
		}
		e.Reply(r,params)
	}


	select {
	case <-time.After(5*time.Second):
		t.Fatal("Didnt Got Result 1")
	case d := <-c1:
		r := d.(*messages.Result)
		if r.Exporter != e.UUID() {
			t.Error("Wrong Export UUID",r.Exporter,e.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params",r.Parameter(),params)
		}
	}
	select {
	case <-time.After(2*time.Second):
		t.Fatal("Didnt Got Result 2")
	case d := <-c2:
		r := d.(*messages.Result)
		if r.Exporter != e.UUID() {
			t.Error("Wrong Export UUID",r.Exporter,e.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params",r.Parameter(),params)
		}
	}


}

func TestCallOne2Many(t *testing.T) {
	i := aurarath.NewImport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer i.Remove()
	i.Run()
	e1 := aurarath.NewExport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer e1.Remove()
	c1 := e1.Requests().AsChan()
	e1.Run()
	e2 := aurarath.NewExport(TEST_APP_DESCRIPTOR,config.DefaultLocalhost())
	defer e2.Remove()
	c2 := e2.Requests().AsChan()
	e2.Run()
	time.Sleep(10*time.Second)
	params := []byte{4,5,63,4}
	params1 := []byte{3}
	params2 := []byte{6}
	s := eventual2go.NewStreamController()
	s1,s2 := s.Split(func(d eventual2go.Data)bool {return d.(*messages.Result).Exporter == e1.UUID()})
	rc1 := s1.AsChan()
	rc2 := s2.AsChan()
	i.CallAll("SayHello",params,s)
	select {
	case <-time.After(5*time.Second):
		t.Fatal("Didnt Got Request 1")
	case d := <-c1:
		r := d.(*messages.Request)
		if r.Importer != i.UUID() {
			t.Error("Wrong Import UUID 1",r.Importer,i.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params 1",r.Parameter(),params)
		}
		e1.Reply(r,params1)
	}


	select {
	case <-time.After(2*time.Second):
		t.Fatal("Didnt Got Request 2")
	case d := <-c2:
		r := d.(*messages.Request)
		if r.Importer != i.UUID() {
			t.Error("Wrong Import UUID 2",r.Importer,i.UUID())
		}
		if !bytes.Equal(r.Parameter(),params){
			t.Error("Wrong Params 2",r.Parameter(),params)
		}
		e2.Reply(r,params2)
	}


	select {
	case <-time.After(2*time.Second):
		t.Fatal("Didnt Got Result 1")
	case d := <-rc1:
		r := d.(*messages.Result)
		if r.Exporter != e1.UUID() {
			t.Error("Wrong Export UUID",r.Exporter,e1.UUID())
		}
		if !bytes.Equal(r.Parameter(),params1){
			t.Error("Wrong Params",r.Parameter(),params1)
		}
	}
	select {
	case <-time.After(2*time.Second):
		t.Fatal("Didnt Got Result 2")
	case d := <-rc2:
		r := d.(*messages.Result)
		if r.Exporter != e2.UUID() {
			t.Error("Wrong Export UUID",r.Exporter,e2.UUID())
		}
		if !bytes.Equal(r.Parameter(),params2){
			t.Error("Wrong Params",r.Parameter(),params2)
		}
	}


}

var TEST_APP_DESCRIPTOR *aurarath.AppDescriptor = &aurarath.AppDescriptor{
	[]aurarath.Function{
		aurarath.Function{
			"SayHello",
			[]aurarath.Parameter{
				aurarath.Parameter{
					"Greeting",
					"string"}},
			[]aurarath.Parameter{
				aurarath.Parameter{
					"Answer",
					"string"}},
		},
	},
	map[string]string{"TAG_1":""},
}