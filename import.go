package aurarath

import "github.com/YamiNoKeshin/aurarath/network"

type Import struct {
	AppKey *AppKey
	nodes  []*network.Node
}

func NewImport(config *Config, key *AppKey) *Import {
	imp := Import{
		AppKey: key,
	}

	for _, iface := range config.NetworkInterfaces {
		append(imp.nodes, network.NewNode(iface, config.Logger()))
	}

	return &imp
}

func (*Import) AddFunction(fkt *Function)

func (*Import) RemoveFunction(name string)

func (*Import) UpdateTags(tags []string)

func (*Import) Call(req *Request) (res *Result)

func (*Import) CallAll(req *Request) (res *Result)

func (*Import) Trigger(req *Request)

func (*Import) TriggerAll(req *Request)

func (*Import) Remove()
