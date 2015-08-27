package service

import (
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"github.com/joernweissenborn/aurarath/config"
	"testing"
	"time"
)

func TestServicBasics(t *testing.T) {
	a := new(appdescriptor.AppDescriptor)
	s1 := NewService(a, EXPORTING, config.DefaultLocalhost(), []byte{0})
	defer s1.Remove()
	c := s1.Disconnected().AsChan()
	c1 := s1.Connected().AsChan()
	s1.Run()
	s2 := NewService(a, IMPORTING, config.DefaultLocalhost(), []byte{0})
	c2 := s2.Connected().AsChan()
	s2.Run()

	select {
	case <-time.After(15 * time.Second):
		t.Error("Service 1 Did Not Connect")

	case <-c1:
	}
	select {
	case <-time.After(1 * time.Second):
		t.Error("Service 2 Did Not Connect")

	case <-c2:
	}

	t.Log("Shutting Down Service 2")

	s2.Remove()
	select {
	case <-time.After(10 * time.Second):
		t.Error("Service 2 Did Not Disconnect")

	case <-c:
		t.Log("Successfully Disconnected Service 2")
	}
}
func TestMultipleServices(t *testing.T) {
	a := new(appdescriptor.AppDescriptor)
	s1 := NewService(a, EXPORTING, config.DefaultLocalhost(), []byte{0})
	defer s1.Remove()
	c1 := s1.Connected().AsChan()
	s1.Run()

	s2 := NewService(a, IMPORTING, config.DefaultLocalhost(), []byte{0})
	defer s2.Remove()
	c2 := s2.Connected().AsChan()
	s2.Run()

	s3 := NewService(a, IMPORTING, config.DefaultLocalhost(), []byte{0})
	defer s3.Remove()
	c3 := s3.Connected().AsChan()
	s3.Run()

	s4 := NewService(a, EXPORTING, config.DefaultLocalhost(), []byte{0})
	defer s4.Remove()
	c4 := s4.Connected().AsChan()
	s4.Run()

	select {
	case <-time.After(15 * time.Second):
		t.Error("Service 1 Did Not Connect")

	case <-c1:
		if len(s1.connectedServices) != 2 {
			t.Error("Service 1 didnt connet to all")
		}
	}
	select {
	case <-time.After(1 * time.Second):
		t.Error("Service 2 Did Not Connect")

	case <-c2:
		if len(s2.connectedServices) != 2 {
			t.Error("Service 2 didnt connet to all")
		}
	}

	select {
	case <-time.After(1 * time.Second):
		t.Error("Service 3 Did Not Connect")

	case <-c3:
		if len(s3.connectedServices) != 2 {
			t.Error("Service 3 didnt connet to all")
		}
	}

	select {
	case <-time.After(1 * time.Second):
		t.Error("Service 4 Did Not Connect")

	case <-c4:
		if len(s3.connectedServices) != 2 {
			t.Error("Service 4 didnt connet to all")
		}
	}

}
