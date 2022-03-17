package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/extensions/go-plugin/pkg/transcoder/bumsbeis"
)

type beis2bums struct {
	cfg       map[string]interface{}
	bums      api.HeaderMap
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
	config, err := bibm.GetConfig()
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
	return bumsHeaders, bumsBuf, trailers, nil
}

func (bibm *beis2bums) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	config, err := bibm.GetConfig()
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
	return beisHeaders, beisBuf, trailers, nil
}

func (bibm *beis2bums) GetConfig() (*Bums2BeisConfig, error) {
	details, ok := bibm.cfg["details"].(string)
	if !ok {
		return nil, fmt.Errorf("the %s of details is not exist", bibm.cfg)
	}
	var configs []*Bums2BeisConfig
	if err := json.Unmarshal([]byte(details), &configs); err != nil {
		return nil, err
	}
	if len(configs) != 1 {
		return nil, fmt.Errorf("the length of configs is illage")
	}
	configs[0].ReqMapping = &bumsbeis.Beis2BumsConfig{
		Path:   configs[0].Path,
		Method: configs[0].Method,
		GWName: configs[0].GWName,
	}
	return configs[0], nil
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
