package main

import (
	"context"

	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/extensions/go-plugin/pkg/transcoder/bumscd"
)

var (
	// 单线程安全
	relations map[string]*config
)

type bums2cd struct {
	cfg      map[string]interface{}
	relation *bumscd.Relation
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	relation, _ := ParseRelation(cfg)
	return &bums2cd{
		cfg:      cfg,
		relation: relation,
	}
}

func (t *bums2cd) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	if t.relation == nil {
		return false
	}
	return true
}

func (t *bums2cd) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	tran, err := bumscd.NewCd2Bums(headers, buf, t.relation)
	if err != nil {
		return headers, buf, trailers, err
	}
	_, err = tran.Body()
	if err != nil {
		return headers, buf, trailers, err
	}
	return headers, buf, trailers, nil
}

func (t *bums2cd) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	_, err := bumscd.NewBums2Beis(headers, buf.String(), t.relation)
	if err != nil {
		return headers, buf, trailers, err
	}
	return headers, buf, trailers, nil
}
