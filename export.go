package aurarath

import "github.com/YamiNoKeshin/aurarath/network"

type Export struct {
	AppKey *AppKey
	nodes  []*network.Node
}

func NewExport(config *Config, key *AppKey) {
	exp := &Export{
		AppKey: key,
	}

	for _, iface := range config.NetworkInterfaces {
		exp.nodes = append(exp.nodes, network.NewNode(iface, config.Logger()))
	}

	return exp
}
