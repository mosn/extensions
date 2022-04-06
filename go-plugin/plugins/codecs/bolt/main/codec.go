package main

import (
	"context"

	"mosn.io/extensions/go-plugin/pkg/protocol/bolt"

	"mosn.io/api"
)

// LoadCodec load codec function
func LoadCodec() api.XProtocolCodec {
	return &Codec{}
}

type Codec struct {
	mapping bolt.StatusMapping
}

func (r Codec) ProtocolName() api.ProtocolName {
	return bolt.ProtocolName
}

func (r Codec) NewXProtocol(context.Context) api.XProtocol {
	return &bolt.BoltProtocol{}
}

func (r Codec) ProtocolMatch() api.ProtocolMatch {
	return bolt.Matcher
}

func (r Codec) HTTPMapping() api.HTTPMapping {
	return r.mapping
}

// compiler check
var _ api.XProtocolCodec = &Codec{}
