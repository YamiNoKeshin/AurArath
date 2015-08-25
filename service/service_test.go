package service
import (
	"testing"
	"time"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/appdescriptor"
)

func TestServicBasics(t *testing.T){
	a :=  new(appdescriptor.AppDescriptor)
	s1 := NewService(a,EXPORTING,config.DefaultLocalhost(),[]byte{0})
	defer s1.Remove()
	c := s1.Disconnected().AsChan()
	c1 := s1.Connected().AsChan()
	s1.Run()
	s2 := NewService(a,IMPORTING,config.DefaultLocalhost(),[]byte{0})
	c2 := s2.Connected().AsChan()
	s2.Run()

select {
case <-time.After(1*time.Second):
		t.Error("Service 1 Did Not Connect")

	case d := <-c1:
		id := d.(string)
		if id != s2.UUID() {
			t.Errorf("Wrong UUID, got %s, want %s",id,s2.UUID())
		}
	}
select {
case <-time.After(1*time.Second):
		t.Error("Service 2 Did Not Connect")

	case d := <-c2:
		id := d.(string)
		if id != s1.UUID() {
			t.Errorf("Wrong UUID, got %s, want %s",id,s1.UUID())
		}
	}


	t.Log("Shutting Down Service 2")

	s2.Remove()
	select {
	case <-time.After(10*time.Second):
		t.Error("Service 2 Did Not Disconnect", s1.disconnected.IsComplete())

	case <-c:
		t.Log("Successfully Disconnected Service 2")
	}
}