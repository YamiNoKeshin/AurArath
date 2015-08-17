package node_test

import (
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/network/node"
	"time"
	"testing"
	"strings"
	"github.com/hashicorp/serf/serf"
)

func TestNodeDiscover(t *testing.T){
	n1 := node.New(config.DefaultLocalhost(),nil)
	defer n1.Shutdown()
	j1 := n1.Join().AsChan()
	n1.Run()

	n2 := node.New(config.DefaultLocalhost(),nil)
	defer n2.Shutdown()
	j2 := n2.Join().AsChan()
	n2.Run()
	t.Log(n1.UUID,n2.UUID)
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("Couldnt find node 1")
	case data := <-j1:
		if !strings.Contains(data.(serf.Member).Name,n2.UUID) {
			t.Error("Found wrong UUID")
		}
	}

	select {
	case <-time.After(5 * time.Minute):
		t.Fatal("Couldnt find node 2")
	case data := <-j2:
		if !strings.Contains(data.(serf.Member).Name,n1.UUID) {
			t.Error("Found wrong UUID")
		}
	}
}

func TestNodeLeave(t *testing.T){
	n1 := node.New(config.DefaultLocalhost(),nil)
	defer n1.Shutdown()
	c := n1.Leave().AsChan()
	n1.Run()

	n2 := node.New(config.DefaultLocalhost(),nil)
	n2.Run()
	time.Sleep(1 *time.Second)
	n2.Leave()
	n2.Shutdown()

	select {
	case <-time.After(5 * time.Second):
		t.Fatal("node didnt leave")
	case data := <-c:
		if !strings.Contains(data.(serf.Member).Name,n2.UUID) {
			t.Error("Found wrong UUID")
		}
	}
}

