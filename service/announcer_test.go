package service

import (
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"testing"
	"time"
)

func TestAnnouncer(t *testing.T) {
	a := new(appdescriptor.AppDescriptor)
	a1 := NewAnnouncer("test1", []string{"127.0.0.1:666"}, EXPORTING, a)
	defer a1.Shutdown()
	a2 := NewAnnouncer("test2", []string{"127.0.0.1:667"}, IMPORTING, a)
	c := a1.ServiceGone().AsChan()
	c1 := a1.ServiceArrived().AsChan()
	c2 := a2.ServiceArrived().AsChan()
	a1.Run()
	a2.Run()

	select {
	case <-time.After(5 * time.Second):
		t.Error("Service 1 Did Not Connect")

	case d := <-c1:
		sa := d.(ServiceArrived)
		if sa.UUID != "test2" {
			t.Errorf("Wrong UUID, got %s, want 'test2'", sa.UUID)
		}
		if sa.Address != "127.0.0.1" {
			t.Errorf("Wrong Addres, got %s, want '127.0.0.1''", sa.Address)
		}
		if sa.Port != 667 {
			t.Errorf("Wrong Port, got %d, want '667'", sa.Port)
		}
	}
	select {
	case <-time.After(1 * time.Second):
		t.Error("Service 1 Did Not Connect")

	case d := <-c2:
		sa := d.(ServiceArrived)
		if sa.UUID != "test1" {
			t.Errorf("Wrong UUID, got %s, want 'test1'", sa.UUID)
		}
		if sa.Address != "127.0.0.1" {
			t.Errorf("Wrong Addres, got %s, want '127.0.0.1''", sa.Address)
		}
		if sa.Port != 666 {
			t.Errorf("Wrong Port, got %d, want '666'", sa.Port)
		}
	}

	t.Log("Shutting Down Service 2")

	a2.Shutdown()
	select {
	case <-time.After(10 * time.Second):
		t.Error("Service 2 Did Not Disconnect")

	case d := <-c:
		sa := d.(ServiceGone)
		if sa.UUID != "test2" {
			t.Errorf("Wrong UUID, got %s, want 'test2'", sa.UUID)
		}
	}
}
