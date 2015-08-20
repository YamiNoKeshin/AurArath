package peer
import (
	"testing"
	"github.com/joernweissenborn/aurarath/network/connection"
	"time"
)

var testip = "127.0.0.1"

func TestPeer(t *testing.T){
	incoming, err := connection.NewIncoming(testip)
	if err != nil {
		t.Fatal(err)
	}
	c := incoming.In().First().AsChan()

	p := New("test")

	p.OpenConnection(testip,incoming.Port())

	err = p.Send([][]byte{[]byte("Hello")})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(100*time.Millisecond):
		t.Error("Didn't get Message")
	case data := <-c:
		if data.([]string)[0] != "test" {
			t.Error("Got Wrong Id")
		}
		if data.([]string)[1] != "Hello" {
			t.Error("Got Wrong Msg")
		}
	}
}