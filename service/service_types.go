package service
import (
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/network/node"
	"log"
)

const (
	EXPORTING = "EXPORT"
	IMPORTING = "IMPORT"
)

func isService(servicetype string)eventual2go.Filter{
	return func(d eventual2go.Data) bool {
		log.Println("ISSERV",servicetype,d)
		return d.(node.QueryEvent).Query.Name == servicetype
	}
}
