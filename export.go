package aurarath

import (
	"fmt"
	"github.com/joernweissenborn/aurarath/appdescriptor"
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/messages"
	"github.com/joernweissenborn/aurarath/service"
	"github.com/joernweissenborn/eventual2go"
	"log"
)

type Export struct {
	*service.Service

	requests *eventual2go.Stream

	r *eventual2go.Reactor

	listeners map[string][]string

	logger *log.Logger
}

func NewExport(a *appdescriptor.AppDescriptor, cfg *config.Config) (e *Export) {
	e = new(Export)
	e.Service = service.NewService(a, service.EXPORTING, cfg, []byte{0})
	e.logger = log.New(cfg.Logger(), fmt.Sprintf("export %s ", e.UUID()), log.Lshortfile)
	e.requests = e.Messages(messages.REQUEST)
	e.listeners = map[string][]string{}
	e.r = eventual2go.NewReactor()
	e.r.React("listen", e.newListener)
	e.r.AddStream("listen", e.IncomingMessages(messages.LISTEN))
	e.r.React("listen_stop", e.stopListener)
	e.r.AddStream("listen", e.IncomingMessages(messages.STOP_LISTEN))
	e.r.React("reply", e.deliverResult)
	return
}

func (e *Export) Reply(r *messages.Request, params []byte) {
	res := messages.NewResult(e.UUID(), r, params)
	e.r.Fire("reply", res)
}

func (e *Export) Requests() *eventual2go.Stream {
	return e.requests
}

func (*Export) Emit(function string, parameter []byte) {}

func (e *Export) newListener(d eventual2go.Data) {
	l := d.(messages.IncomingMessage)
	f := l.Msg.(*messages.Listen).Function
	e.logger.Println("New Listener", l.Sender, f)
	ls := e.listeners[f]
	if ls == nil {
		e.listeners[f] = []string{l.Sender}
		return
	}
	for _, id := range ls {
		if id == l.Sender {
			return
		}
	}
	e.listeners[f] = append(ls, l.Sender)
}

func (e *Export) removeListener(d eventual2go.Data) {
	l := d.(string)
	for f, ls := range e.listeners {
		index := -1

		for i, id := range ls {
			if id == l {
				index = i
				break
			}
		}
		if index == -1 {
			continue
		}
		ls[index] = ls[len(ls)-1]
		e.listeners[f] = ls[:len(ls)-2]
	}
}

func (e *Export) stopListener(d eventual2go.Data) {
	l := d.(messages.IncomingMessage)
	f := l.Msg.(*messages.StopListen).Function
	ls := e.listeners[f]
	if ls == nil {
		return
	}
	index := -1
	for i, id := range ls {
		if id == l.Sender {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	ls[index] = ls[len(ls)-1]
	e.listeners[f] = ls[:len(ls)-2]
}

func (e *Export) deliverResult(d eventual2go.Data) {
	result := d.(*messages.Result)
	e.logger.Println("Delivering result", result.Request.Function, result.Request.CallType)
	switch result.Request.CallType {
	case messages.ONE2MANY, messages.ONE2ONE:
		if sc := e.GetConnectedService(result.Request.Importer); sc != nil {
			sc.Send(messages.Flatten(result))
		}

	case messages.MANY2MANY, messages.MANY2ONE:
		res := messages.Flatten(result)
		e.logger.Printf("sending many2 result to %d clients", len(e.listeners))
		for _, uuid := range e.listeners[result.Request.Function] {
			e.logger.Println("Sending result", uuid, result.Request.Function)
			e.GetConnectedService(uuid).Send(res)
		}
	}
}
