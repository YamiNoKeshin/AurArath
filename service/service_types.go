package service

import (
	"github.com/joernweissenborn/aurarath/network/node"
	"github.com/joernweissenborn/eventual2go"
)

const (
	EXPORTING = "EXPORT"
	IMPORTING = "IMPORT"
)

func isService(servicetype string) eventual2go.Filter {
	return func(d eventual2go.Data) bool {
		return d.(node.QueryEvent).Query.Name == servicetype
	}
}
