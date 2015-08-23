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
	e.listeners = map[string][]string{}
	e.in.Where(messages.Is(messages.LISTEN)).Listen(e.newListener)
	e.in.Where(messages.Is(messages.STOP_LISTEN)).Listen(e.stopListener)
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

func (e *Export) newListener(d eventual2go.Data){
	e.m.Lock()
	defer e.m.Unlock()
	l := d.(messages.IncomingMessage)
	f := l.Msg.(*messages.Listen).Function
	e.logger.Println("New Listener",l.Sender,f)
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
	e.listeners[f] = append(ls,l.Sender)
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
		ls[index] = ls[len(ls) - 1]
		e.listeners[f] = ls[:len(ls) - 2]
	}
}
func (e *Export) stopListener(d eventual2go.Data){
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

func (e *Export) deliverResult(result *messages.Result){
	e.m.Lock()
	defer e.m.Unlock()
	e.logger.Println("Delivering result",result.Request.Function,result.Request.CallType)
	switch result.Request.CallType {
	case messages.ONE2MANY, messages.ONE2ONE:
		if p := e.getPeer(result.Request.Importer); p != nil {
			p.Send(messages.Flatten(result))
		}

	case messages.MANY2MANY, messages.MANY2ONE:
		res := messages.Flatten(result)
		e.logger.Println("Delivering MANY2",result.Request.Function)
		e.logger.Println("Delivering MANY2",e.listeners)
		for k := range e.listeners {
			e.logger.Println("Delivering MANY2",k)
			e.logger.Println("Delivering MANY2",[]byte(k))
			e.logger.Println("Delivering MANY2",[]byte(result.Request.Function))

		}
		e.logger.Println("Delivering MANY2",e.listeners[result.Request.Function])
		for _, pid := range e.listeners[result.Request.Function] {
			e.logger.Println("Sending res",pid,result.Request.Function)
			e.getPeer(pid).Send(res)
		}
	}
}