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

type Import struct {
	*service.Service

	pending []*messages.Request

	results *eventual2go.Stream

	listen []string

	logger *log.Logger
}

func NewImport(a *appdescriptor.AppDescriptor, cfg *config.Config) (i *Import) {
	i = new(Import)
	i.Service = service.NewService(a, service.IMPORTING, cfg, []byte{0})
	i.logger = log.New(cfg.Logger(), fmt.Sprintf("export %s ", i.UUID()), log.Lshortfile)
	i.Connected().Then(i.onConnected)
	i.NewServiceConnections().Listen(i.sendListenFunctions)
	i.results = i.Messages(messages.RESULT)
	return
}

func (i *Import) sendListenFunctions(d eventual2go.Data) {
	sc := d.(*service.ServiceConnection)
	for _, f := range i.listen {
		sc.Send(messages.Flatten(&messages.Listen{f}))
	}
	return
}

func (i *Import) Call(function string, parameter []byte) (f *eventual2go.Future) {
	i.logger.Println("Call", function)
	uuid := i.call(function, parameter, messages.ONE2ONE)
	f = i.results.FirstWhere(isRes(uuid))
	return
}

func (i *Import) CallAll(function string, parameter []byte, s *eventual2go.StreamController) {
	i.logger.Println("CallAll", function)
	uuid := i.call(function, parameter, messages.ONE2MANY)
	s.Join(i.results.Where(isRes(uuid)))
	return
}

func (i *Import) Listen(function string) {
	for _, f := range i.listen {
		if f == function {
			return
		}
	}
	i.listen = append(i.listen, function)
	for _, sc := range i.GetConnectedServices() {
		sc.Send(messages.Flatten(&messages.Listen{function}))
	}
}

func (i *Import) StopListen(function string) {
	index := -1
	for i, f := range i.listen {
		if f == function {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	i.listen[index] = i.listen[len(i.listen)-1]
	i.listen = i.listen[:len(i.listen)-2]
	for _, sc := range i.GetConnectedServices() {
		sc.Send(messages.Flatten(&messages.StopListen{function}))
	}
}

func (i *Import) Trigger(function string, parameter []byte) {
	i.call(function, parameter, messages.MANY2ONE)
}

func (i *Import) TriggerAll(function string, parameter []byte) {
	i.call(function, parameter, messages.MANY2MANY)
}

func (i *Import) Results() *eventual2go.Stream {
	return i.results.Where(func(d eventual2go.Data) bool {
		return d.(*messages.Result).Request.CallType == messages.MANY2MANY || d.(*messages.Result).Request.CallType == messages.MANY2ONE
	})
}

func (i *Import) call(function string, parameter []byte, ctype messages.CallType) (uuid string) {

	req := messages.NewRequest(i.UUID(), function, ctype, parameter)
	if i.Service.Connected().Completed() {
		i.deliverRequest(req)
	} else {
		i.pending = append(i.pending, req)
	}
	return req.UUID
}

func isRes(uuid string) eventual2go.Filter {
	return func(d eventual2go.Data) bool {
		return d.(*messages.Result).Request.UUID == uuid
	}
}

func (i *Import) deliverRequest(r *messages.Request) {
	for _, sc := range i.GetConnectedServices() {
		i.logger.Println("Delivering Request to", sc.Uuid())
		sc.Send(messages.Flatten(r))
		if r.CallType == messages.ONE2ONE || r.CallType == messages.MANY2ONE {
			return
		}
	}
	return
}

func (i *Import) onConnected(eventual2go.Data) eventual2go.Data {
	i.deliverAllRequests()
	return nil
}

func (i *Import) deliverAllRequests() {
	for _, r := range i.pending {
		i.deliverRequest(r)
	}
}
