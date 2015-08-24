package daemon
import (
	"strconv"
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/network/connection"
	"bytes"
	"github.com/joernweissenborn/aurarath"
	"encoding/json"
)

const PROTOCOL_SIGNATURE byte = 0xA1

const (
	HELLO int = iota
	EXPORT
	IMPORT
	MEMBER_FOUND
	MEMBER_GONE
)

func NewHello(port int) [][]byte {
	return [][]byte{PROTOCOL_SIGNATURE,HELLO,[]byte(strconv.FormatInt(port,10))}
}

func IsMessage(t int) eventual2go.Filter {
	return func(d eventual2go.Data) bool {
		m := d.(connection.Message).Payload
		return  m[2] == t
	}
}

func ValidMessage(d eventual2go.Data) bool {
	m := d.(connection.Message).Payload
	if !bytes.Equal(m[1],PROTOCOL_SIGNATURE) {
		return false
	}
	return m[2] == EXPORT,m[2] == IMPORT
}

type NewService struct {
	UUID string
	Descriptor *aurarath.AppDescriptor
	Addresses []string
	ServiceType string
}

func ToNewServiceMessage(d eventual2go.Data) eventual2go.Data{
	m := d.(connection.Message).Payload
	var nsm NewService
	json.Unmarshal(m[3],&nsm)
	return nsm
}

type MemberFound struct {
	UUID string
	Address string
	Port int
}

func NewMemberFound(uuid, address string, port int) [][]byte {
	b, _ := json.Marshal(MemberFound{uuid,address,port})
	return [][]byte{[]byte{PROTOCOL_SIGNATURE},[]byte{MEMBER_FOUND},b}
}
type MemberGone struct {
	UUID string
}

func NewMemberGone(uuid string) [][]byte {
	b, _ := json.Marshal(MemberGone{uuid})
	return [][]byte{[]byte{PROTOCOL_SIGNATURE},[]byte{MEMBER_GONE},b}
}