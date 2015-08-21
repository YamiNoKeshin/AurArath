package messages

//go:generate stringer -type=Codec

type Codec int

const (
	JSON CallType = iota
	MSGPACK
	BSON
)