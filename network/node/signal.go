package node

import (
	"encoding/binary"
	"github.com/joernweissenborn/aurarath/network/beacon"
	"github.com/joernweissenborn/eventual2go"
	"net"
)

const PROTOCOLL_SIGNATURE uint8 = 0xA5

type PeerAddress struct {
	IP   net.IP
	Port uint16
}

func SignalToAdress(d eventual2go.Data) eventual2go.Data {
	var a PeerAddress
	a.IP = d.(beacon.Signal).SenderIp
	a.Port = binary.LittleEndian.Uint16(d.(beacon.Signal).Data[1:])
	return a
}

func IsValidSignal(d eventual2go.Data) bool {
	if len(d.(beacon.Signal).Data) == 0 {
		return false
	}
	sig := d.(beacon.Signal).Data[0]
	return sig == PROTOCOLL_SIGNATURE
}

func NewSignalPayload(port uint16) (payload []byte) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, port)
	payload = []byte{PROTOCOLL_SIGNATURE, b[0], b[1]}
	return
}
