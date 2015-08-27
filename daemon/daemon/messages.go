package daemon

import (
	"encoding/json"
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"github.com/joernweissenborn/aurarath/network/connection"
	"github.com/joernweissenborn/aurarath/service"
	"github.com/joernweissenborn/eventual2go"
	"strconv"
)

const PROTOCOL_SIGNATURE byte = 0xA1

const (
	HELLO uint8 = iota
	EXPORT
	IMPORT
	SERVICE_ARRIVE
	SERVICE_GONE
)

func NewHello(port int) [][]byte {
	return [][]byte{[]byte{byte(PROTOCOL_SIGNATURE)}, []byte{byte(HELLO)}, []byte(strconv.FormatInt(int64(port), 10))}
}

func IsMessage(t uint8) eventual2go.Filter {
	return func(d eventual2go.Data) bool {
		m := d.(connection.Message).Payload
		return []byte(m[2])[0] == t
	}
}

func ValidMessage(d eventual2go.Data) bool {
	m := d.(connection.Message).Payload
	return m[1][0] == PROTOCOL_SIGNATURE
}

type NewService struct {
	UUID        string
	Descriptor  *appdescriptor.AppDescriptor
	Addresses   []string
	ServiceType string
}

func ToNewServiceMessage(d eventual2go.Data) eventual2go.Data {
	m := d.(connection.Message).Payload
	var nsm NewService
	json.Unmarshal([]byte(m[3]), &nsm)
	return nsm
}

func ToServiceArrivedMessage(d eventual2go.Data) eventual2go.Data {
	m := d.(connection.Message).Payload
	var nsm service.ServiceArrived
	json.Unmarshal([]byte(m[3]), &nsm)
	return nsm
}

func ToServiceGone(d eventual2go.Data) eventual2go.Data {
	m := d.(connection.Message).Payload
	return m[3]
}

func (ns NewService) flatten() [][]byte {
	b, _ := json.Marshal(ns)

	var m uint8

	if ns.ServiceType == service.EXPORTING {
		m = EXPORT
	} else {
		m = IMPORT
	}

	return [][]byte{[]byte{PROTOCOL_SIGNATURE}, []byte{m}, b}
}

func NewServiceArrived(sa service.ServiceArrived) [][]byte {
	b, _ := json.Marshal(sa)
	return [][]byte{[]byte{PROTOCOL_SIGNATURE}, []byte{SERVICE_ARRIVE}, b}
}

func NewServiceGone(uuid string) [][]byte {
	return [][]byte{[]byte{PROTOCOL_SIGNATURE}, []byte{SERVICE_GONE}, []byte(uuid)}
}
