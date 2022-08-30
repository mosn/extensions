package generator

import (
	"context"

	"mosn.io/api"
)

type ProtocolDelegate func(ctx context.Context, frame api.XFrame, span api.Span)

var (
	delegates = make(map[api.ProtocolName]ProtocolDelegate)
)

func RegisterDelegate(name api.ProtocolName, delegate ProtocolDelegate) {
	delegates[name] = delegate
}

func GetDelegate(name api.ProtocolName) ProtocolDelegate {
	return delegates[name]
}
