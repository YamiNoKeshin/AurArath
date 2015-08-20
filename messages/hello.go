package messages
import (
	"github.com/joernweissenborn/aurarath/network/peer"
)


type Hello struct {
	PeerDetails peer.Details
}

func(Hello) GetType() uint16 {return HELLO}
