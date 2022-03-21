package main

import (
	"context"

	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/extensions/go-plugin/pkg/transcoder/bumscd"
)

type cd2bums struct {
	cfg    map[string]interface{}
	config *bumscd.Config
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	return &cd2bums{
		cfg: cfg,
	}
}

func (t *cd2bums) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	config, err := t.ParseConfig(t.cfg)
	if err != nil {
		return false
	}
	t.config = config
	return true
}

func (t *cd2bums) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	tran, err := bumscd.NewCd2Bums(ctx, headers, buf, t.config)
	if err != nil {
		return headers, buf, trailers, err
	}
	cdHeaders, cdBuf, err := tran.Transcoder(false)
	if err != nil {
		return headers, buf, trailers, err
	}
	return cdHeaders, cdBuf, trailers, nil
}

func (t *cd2bums) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	cd2bums, err := bumscd.NewBums2Cd(ctx, headers, buf, t.config)
	if err != nil {
		return headers, buf, trailers, err
	}
	bumsHeaders, bumsBuf, err := cd2bums.Transcoder(true)
	if err != nil {
		return headers, buf, trailers, nil
	}
	return bumsHeaders, bumsBuf, trailers, nil
}

func (t *cd2bums) ParseConfig(cfg map[string]interface{}) (*bumscd.Config, error) {
	rInfo, ok := cfg["details"]
	if !ok {
		return nil, nil
	}
	info, ok := rInfo.(string)
	if !ok {
		return nil, nil
	}
	return configManager.GetLatestRelation(info)
}
