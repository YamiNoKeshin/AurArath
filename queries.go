package connection
import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/network/node"
)

const (
	EXPORTING = "EXPORT"
	IMPORTING = "IMPORT"
)

func isService(servicetype string)eventual2go.Filter{
	return func(d eventual2go.Data) bool {
		return d.(node.QueryEvent).Query.Name == servicetype
	}
}
