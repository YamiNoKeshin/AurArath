package messages
import "encoding/gob"

const (
	HELLO = iota
)

func init(){
	gob.Register(Hello{})
}