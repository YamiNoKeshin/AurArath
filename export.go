package connection

import (
	"aurarath/config"
	"github.com/joernweissenborn/aurarath/network/node"
)


type Export struct {
	appDescriptor *AppDescriptor

	node   *node.Node
}

func NewExport(a *AppDescriptor, cfg *config.Config) (e *Export){
	e = new(Export)
	e.appDescriptor = a
	e.node = node.New(cfg,a.AsTagSet())
	e.node.Queries().Where(IsExporting)
	return
}



