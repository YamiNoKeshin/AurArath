package daemon_test
import "testing"
import (
	"github.com/joernweissenborn/aurarath/daemon/daemon"
	"github.com/joernweissenborn/aurarath/service"
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"time"
)

func TestDaemon(t *testing.T) {
	daemon.New("127.0.0.1",5558)
	ad := new(appdescriptor.AppDescriptor)

	c1 := daemon.NewClient("TESTID1","127.0.0.1","127.0.0.1",5558,service.EXPORTING,ad,[]string{"127.0.0.1:666"})
	arrive1 := c1.Arrived().AsChan()
	c1.Run()
	c2 := daemon.NewClient("TESTID2","127.0.0.1","127.0.0.1",5558,service.IMPORTING,ad,[]string{"127.0.0.1:667"})
	arrive2 := c2.Arrived().AsChan()
	c2.Run()
	select {
	case <-time.After(20*time.Second):
		t.Fatal("import didnt arrive")
	case d := <-arrive1:
		r := d.(service.ServiceArrived)
		if r.UUID != "TESTID2"{
			t.Error("Wrong UUID",r.UUID,"TESTID2")
		}
		if r.Interface != "127.0.0.1"{
			t.Error("Wrong interface",r.Interface, "127.0.0.1")
		}
		if r.Port != 667{
			t.Error("Wrong port",r.Port, 667)
		}
	}

	select {
	case <-time.After(2*time.Second):
		t.Fatal("export didnt arrive")
	case d := <-arrive2:
		r := d.(service.ServiceArrived)
		if r.UUID != "TESTID1"{
			t.Error("Wrong UUID",r.UUID,"TESTID1")
		}
		if r.Interface != "127.0.0.1"{
			t.Error("Wrong interface",r.Interface, "127.0.0.1")
		}
		if r.Port != 666{
			t.Error("Wrong port",r.Port, 666)
		}
	}

}