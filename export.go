package aurarath

import (
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/aurarath/messages"
)


type Export struct {
	*Service

	requests eventual2go.Stream

	listeners map[string][]string

}

func NewExport(a *AppDescriptor, cfg *config.Config) (e *Export){
	e = new(Export)
	e.Service = NewService(a, EXPORTING, cfg,[]byte{0})
	e.requests = e.in.Where(messages.Is(messages.REQUEST)).Transform(messages.ToMsg)
	return
}

func (e *Export) Reply(r *messages.Request, params []byte){
	res := messages.NewResult(e.UUID(),r,params)
	e.deliverResult(res)
}

func (e *Export) Requests() eventual2go.Stream {
	return e.requests
}

func (* Export) Emit(function string, parameter []byte){}

func (e *Export) deliverResult(result *messages.Result){
	switch result.Request.CallType {
	case messages.ONE2MANY, messages.ONE2ONE:
		if p := e.getPeer(result.Request.Importer); p != nil {
			p.Send(messages.Flatten(result))
		}

	case messages.MANY2MANY, messages.MANY2ONE:
		res := messages.Flatten(result)
		for _, pid := range e.listeners[result.Request.Function] {
			e.getPeer(pid).Send(res)
		}
	}
}