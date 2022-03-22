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

type beis2bums struct {
	cfg       map[string]interface{}
	namespace string
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	return &beis2bums{
		cfg: cfg,
	}
}

func (bibm *beis2bums) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	if err := bibm.PraseeeNamespace(headers); err != nil {
		return false
	}
	//TODO 参数配置解析
	return true
}

func (bibm *beis2bums) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	config, err := bibm.GetConfig(ctx)
	if err != nil {
		return headers, buf, trailers, err
	}

	vo := &bumsbeis.Bums2BeisVo{
		Namespace: bibm.namespace,
		GWName:    config.GWName,
	}
	br2br, err := bumsbeis.NewBums2Beis(ctx, headers, buf, config.RespMapping, vo)
	if err != nil {
		return headers, buf, trailers, nil
	}

	if err := br2br.CheckParam(); err != nil {
		return headers, buf, trailers, err
	}

	beisHeaders, beisBuf, err := br2br.Transcoder(false)
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

func (bibm *beis2bums) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	config, err := bibm.GetConfig(ctx)
	if err != nil {
		return headers, buf, trailers, err
	}
	br2br, err := bumsbeis.NewBeis2Bums(ctx, headers, buf, config.ReqMapping)
	if err != nil {
		return headers, buf, trailers, nil
	}
	bumsHeaders, bumsBuf, err := br2br.Transcoder(true)
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

func (bibm *beis2bums) GetConfig(ctx context.Context) (*Bums2BeisConfig, error) {
	details, ok := bibm.cfg["details"].(string)
	if !ok {
		return nil, fmt.Errorf("the %s of details is not exist", bibm.cfg)
	}
	cfg, err := configManager.GetLatestRelation(details)
	if err != nil {
		return nil, err
	}
	if log.DefaultContextLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultContextLogger.Debugf(ctx, "[transcoders][beis2bums] config:%s", details)
	}
	return cfg, nil
}

func (bibm *beis2bums) PraseeeNamespace(headers api.HeaderMap) error {
	scence, ok := headers.Get("ServiceScene")
	if !ok {
		return fmt.Errorf("the key of scence is not exist")
	}
	service, ok := headers.Get("ServiceCode")
	if !ok {
		return fmt.Errorf("the key of service is not exist")
	}
	bibm.namespace = strings.ToLower(scence) + "." + service
	return nil
}
