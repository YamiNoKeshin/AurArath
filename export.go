package connection

import (
	"github.com/joernweissenborn/aurarath/config"
)


type Export struct {
	*Service
}

func NewExport(a *AppDescriptor, cfg *config.Config) (e *Export){
	e = new(Export)
	e.Service = NewService(a, EXPORTING, cfg,[]byte{0})
	return
}




