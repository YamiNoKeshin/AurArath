package connection
import (
	"testing"
	"time"
	"github.com/joernweissenborn/aurarath/config"
)

func TestServiceFindingBasics(t *testing.T){
	a :=  new(AppDescriptor)
	s1 := NewService(a,EXPORTING,config.DefaultLocalhost(),[]byte{0})
	s2 := NewService(a,IMPORTING,config.DefaultLocalhost(),[]byte{0})
	s1.Run()
	s2.Run()
	time.Sleep(1*time.Second)
	if len(s1.peers) == 0 {
		t.Error("Didn't found service 2")
	}
	if len(s2.peers) == 0 {
		t.Error("Didn't found service 1")
	}
}