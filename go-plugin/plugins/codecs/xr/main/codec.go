package main

import (
	"context"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/xr"
	"mosn.io/api"
)

// LoadCodec load codec function
func LoadCodec() api.XProtocolCodec {
	return &Codec{}
}

type Codec struct {
	proto               xr.XrProtocol
	xrMatcher           xr.Matcher
	xrHttpStatusMapping xr.StatusMapping
}

func (r Codec) ProtocolName() api.ProtocolName {
	return r.proto.Name()
}

func (r Codec) ProtocolMatch() api.ProtocolMatch {
	return r.xrMatcher.XrProtocolMatcher
}

func (r Codec) HTTPMapping() api.HTTPMapping {
	return &r.xrHttpStatusMapping
}

func (r Codec) NewXProtocol(ctx context.Context) api.XProtocol {
	return &xr.XrProtocol{}
}

// compiler check
var _ api.XProtocolCodec = &Codec{}
