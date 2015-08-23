package aurarath

import (
	"github.com/joernweissenborn/aurarath/config"
	"github.com/joernweissenborn/aurarath/messages"
	"github.com/joernweissenborn/eventual2go"
)

type Import struct {
	*Service

	pending []*messages.Request

	results eventual2go.Stream
}

func NewImport(a *AppDescriptor, cfg *config.Config) (i *Import){
	i = new(Import)
	i.Service = NewService(a, IMPORTING, cfg,[]byte{0})
	i.connected.Then(i.onConnected)
	i.results = i.in.Where(messages.Is(messages.RESULT)).Transform(messages.ToMsg)
	return
}

func (*Import) UpdateTags(tags []string)   {}

func (i *Import) Call(function string, parameter []byte) (f *eventual2go.Future)     {
	uuid := i.call(function,parameter, messages.ONE2ONE)
	f = i.results.FirstWhere(isRes(uuid))
	return
}

func (*Import) CallAll()     {
	return
}

func (*Import) Trigger() {}

func (*Import) TriggerAll()   {}


func (i *Import) call(function string, parameter []byte, ctype messages.CallType) (uuid string){

	req := messages.NewRequest(i.UUID(),function,ctype,parameter)
	if i.Service.connected.IsComplete() {
		i.deliverRequest(req)
	} else {
		i.pending = append(i.pending,req)
	}
	return req.UUID
}

func isRes(uuid string) eventual2go.Filter {
	return func (d eventual2go.Data) bool {
		return d.(*messages.Result).Request.UUID == uuid
	}
}




func (i *Import) deliverRequest(r *messages.Request) {
	for _, p := range i.getConnectedPeers() {
		p.Send(messages.Flatten(r))
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
