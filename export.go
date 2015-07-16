package aurarath

import "github.com/YamiNoKeshin/aurarath/network"

type Export struct {
	AppKey *AppKey
	node   *network.Node
}

func NewExport(config *Config, key *AppKey) {
	exp := Export{
		AppKey: key,
		node:   network.NewNode(config),
	}

	return &exp
}
