package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/extensions/go-plugin/pkg/transcoder/bumsbeis"
	"mosn.io/pkg/log"
)

type bums2beis struct {
	cfg  map[string]interface{}
	bums api.HeaderMap
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	return &bums2beis{
		cfg: cfg,
	}
}

func (bmbi *bums2beis) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	//TODO 参数配置解析
	return true
}

func (bmbi *bums2beis) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	config, vo, err := bmbi.GetConfig(ctx)
	if err != nil {
		return headers, buf, trailers, err
	}
	br2br, err := bumsbeis.NewBums2Beis(ctx, headers, buf, config, vo)
	if err != nil {
		return headers, buf, trailers, nil
	}

	if err := br2br.CheckParam(); err != nil {
		return headers, buf, trailers, err
	}

	beisHeaders, beisBuf, err := br2br.Transcoder(true)
	if err != nil {
		return headers, buf, trailers, err
	}

	if log.DefaultContextLogger.GetLogLevel() >= log.DEBUG {
		jhs, _ := json.Marshal(headers)
		jhd, _ := json.Marshal(beisHeaders)
		log.DefaultContextLogger.Debugf(ctx, "[transcoders][beis2bums] tran request src_head:%s,dst_head:%s,src_body:%s,dst_body:%s", jhs, jhd, buf.String(), beisBuf.String())
	}
	return beisHeaders, beisBuf, trailers, nil
}

func (bmbi *bums2beis) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	br2br, err := bumsbeis.NewBeis2Bums(ctx, headers, buf, nil)
	if err != nil {
		return headers, buf, trailers, nil
	}
	bumsHeaders, bumsBuf, err := br2br.Transcoder(false)
	if err != nil {
		return headers, buf, trailers, nil
	}
	if log.DefaultContextLogger.GetLogLevel() >= log.DEBUG {
		jhs, _ := json.Marshal(headers)
		jhd, _ := json.Marshal(bumsHeaders)
		log.DefaultContextLogger.Debugf(ctx, "[transcoders][beis2bums] tran request src_head:%s,dst_head:%s,src_body:%s,dst_body:%s", jhs, jhd, buf.String(), bumsBuf.String())
	}
	return bumsHeaders, bumsBuf, trailers, nil
}

func (bmbi *bums2beis) GetConfig(ctx context.Context) (*bumsbeis.Bums2BeisConfig, *bumsbeis.Bums2BeisVo, error) {
	details, ok := bmbi.cfg["details"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("the %s of details is not exist", bmbi.cfg)
	}
	cfg, err := configManager.GetLatestRelation(details)
	if err != nil {
		return nil, nil, err
	}
	vo := &bumsbeis.Bums2BeisVo{
		Namespace: strings.ToLower(cfg.ServiceScene) + "." + cfg.ServiceCode,
		GWName:    cfg.GWName,
	}
	if log.DefaultContextLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultContextLogger.Debugf(ctx, "[transcoders][beis2bums] config:%s", details)
	}
	return cfg.ReqMapping, vo, nil
}
