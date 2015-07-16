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
