package main

import (
	"context"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/cd"
	"mosn.io/api"
)

// LoadCodec load codec function
func LoadCodec() api.XProtocolCodec {
	return &Codec{}
}

type Codec struct {
	proto               cd.Protocol
	xrMatcher           cd.Matcher
	xrHttpStatusMapping cd.StatusMapping
}

func (r Codec) ProtocolName() api.ProtocolName {
	return r.proto.Name()
}

func (r Codec) ProtocolMatch() api.ProtocolMatch {
	return r.xrMatcher.CdProtocolMatcher
}

func (r Codec) HTTPMapping() api.HTTPMapping {
	return &r.xrHttpStatusMapping
}

func (r Codec) NewXProtocol(ctx context.Context) api.XProtocol {
	return &cd.Protocol{}
}

// compiler check
var _ api.XProtocolCodec = &Codec{}
