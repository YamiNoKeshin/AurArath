package aurarath

import (
	"github.com/joernweissenborn/aurarath/network/node"
	"github.com/joernweissenborn/aurarath/config"
)

type Import struct {
	*Service
}

func NewImport(a *AppDescriptor, cfg *config.Config) (e *Export){
	e = new(Export)
	e.Service = NewService(a, IMPORTING, cfg,[]byte{0})
	return
}



func (*Import) AddFunction(fkt *Function)  {
}

func (*Import) RemoveFunction(name string)   {}

func (*Import) UpdateTags(tags []string)   {}

func (*Import) Call(req *Request) (res *Result)     {
	return nil
}

func (*Import) CallAll(req *Request) (res *Result)    {
	return nil
}

func (*Import) Trigger(req *Request)          {}

func (*Import) TriggerAll(req *Request)   {}

func (*Import) Remove(){}
