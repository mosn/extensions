package main

import (
	"context"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
)

// LoadCodec load codec function
func LoadCodec() api.XProtocolCodec {
	return &Codec{}
}

type Codec struct {
	proto                 beis.Protocol
	beisMatcher           beis.Matcher
	beisHttpStatusMapping beis.StatusMapping
}

func (r Codec) ProtocolName() api.ProtocolName {
	return r.proto.Name()
}

func (r Codec) ProtocolMatch() api.ProtocolMatch {
	return r.beisMatcher.BeisProtocolMatcher
}

func (r Codec) HTTPMapping() api.HTTPMapping {
	return &r.beisHttpStatusMapping
}

func (r Codec) NewXProtocol(ctx context.Context) api.XProtocol {
	return &beis.Protocol{}
}

// compiler check
var _ api.XProtocolCodec = &Codec{}
