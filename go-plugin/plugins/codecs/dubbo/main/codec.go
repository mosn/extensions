package main

import (
	"context"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo"
	"mosn.io/api"
)

// LoadCodec load codec function
func LoadCodec() api.XProtocolCodec {
	return &Codec{}
}

type Codec struct {
	HttpStatusMapping dubbo.StatusMapping
}

func (r Codec) ProtocolName() api.ProtocolName {
	return dubbo.ProtocolName
}

func (r Codec) ProtocolMatch() api.ProtocolMatch {
	return dubbo.Matcher
}

func (r Codec) HTTPMapping() api.HTTPMapping {
	return r.HttpStatusMapping
}

func (r Codec) NewXProtocol(context.Context) api.XProtocol {
	return dubbo.DubboProtocol{}
}

// compiler check
var _ api.XProtocolCodec = &Codec{}
