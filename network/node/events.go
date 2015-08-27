package node

import (
	"github.com/hashicorp/serf/serf"
	"github.com/joernweissenborn/eventual2go"
)

type QueryEvent struct {
	Address string
	Query   *serf.Query
}

func toQueryEvent(iface string) eventual2go.Transformer {
	return func(d eventual2go.Data) eventual2go.Data {
		return QueryEvent{iface, d.(*serf.Query)}
	}
}

type QueryResponseEvent struct {
	Interface string
	Response  serf.NodeResponse
}

func toQueryResponseEvent(iface string) eventual2go.Transformer {
	return func(d eventual2go.Data) eventual2go.Data {
		return QueryResponseEvent{iface, d.(serf.NodeResponse)}
	}
}

type LeaveEvent struct {
	Name string
}

func toLeaveEvent() eventual2go.Transformer {
	return func(d eventual2go.Data) eventual2go.Data {
		return LeaveEvent{d.(serf.Member).Name}
	}
}
