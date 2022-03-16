package main

import (
	"context"

	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
)

type bums2beis struct {
	cfg map[string]interface{}
	//

	MesgId int64
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	return &bums2beis{
		cfg: cfg,
	}
}

func (t *bums2beis) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

func (t *bums2beis) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	// TODO vo &
	// vo      Bums2BeisVo
	config := Bums2BeisConfig{}
	br2br, err := NewBumsReq2BeisReq(headers, "", config)
	if err != nil {
		return headers, buf, trailers, nil
	}
	if br2br.CheckParam() {
		return headers, buf, trailers, nil
	}
	beisHeaders, beisBuf, err := br2br.Transcoder()
	if err != nil {
		return headers, buf, trailers, nil
	}
	beisHeaders.Set("Content-Type", "application/json")
	return beisHeaders, beisBuf, trailers, nil
}

func (t *bums2beis) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	br2br, err := NewBeisResp2BumsResp(headers, buf)
	if err != nil {
		return headers, buf, trailers, nil
	}
	bumsHeaders, bumsBuf, err := br2br.Transcoder()
	if err != nil {
		return headers, buf, trailers, nil
	}
	// TODO ID
	return bumsHeaders, bumsBuf, trailers, nil
}
