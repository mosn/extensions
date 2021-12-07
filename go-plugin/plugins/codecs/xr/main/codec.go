package main

import (
	"github.com/mosn/extensions/go-plugin/pkg/protocol/xr"
	"mosn.io/api"
)

// LoadCodec load codec function
func LoadCodec() api.XProtocolCodec {
	return &Codec{}
}

type Codec struct {
	proto               xr.Proto
	xrMatcher           xr.Matcher
	xrHttpStatusMapping xr.StatusMapping
}

func (r Codec) ProtocolName() api.ProtocolName {
	return r.proto.Name()
}

func (r Codec) XProtocol() api.XProtocol {
	return &r.proto
}

func (r Codec) ProtocolMatch() api.ProtocolMatch {
	return r.xrMatcher.XrProtocolMatcher
}

func (r Codec) HTTPMapping() api.HTTPMapping {
	return &r.xrHttpStatusMapping
}

// NewProtocolFactory create protocol per stream connection
func (r Codec) NewProtocolFactory() api.XProtocolFactory {
	return &r
}

func (r Codec) NewXProtocol() api.XProtocol {
	return &xr.Proto{}
}

// compiler check
var _ api.XProtocolFactory = &Codec{}
var _ api.XProtocolCodec = &Codec{}
