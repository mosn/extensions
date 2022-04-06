package main

import (
	"context"
	"encoding/json"
	"fmt"

	"mosn.io/api"
)

type Config struct {
	UniqueId    string      `json:"unique_id"`
	Path        string      `json:"path"`
	Method      string      `json:"method"`
	TragetApp   string      `json:"target_app"`
	ReqMapping  ReqMapping  `json:"req_mapping"`
	RespMapping interface{} `json:"-"`
}

type ReqMapping struct {
	Query      []*query         `json:"query"`
	Body       *body            `json:"body"`
	PathParams []*httpPathParam `json:"path_params"`
}

type query struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type body struct {
	Type string `json:"type"`
}

type httpPathParam struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

func (t *dubbo2springcloud) getConfig(ctx context.Context, headers api.HeaderMap) (*Config, error) {
	details, ok := t.cfg["details"]
	if !ok {
		return nil, fmt.Errorf("the %s of details is not exist", t.cfg)
	}

	binfo, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	var cfgs []*Config
	if err := json.Unmarshal(binfo, &cfgs); err != nil {
		return nil, err
	}
	method, _ := headers.Get("method")
	for _, cfg := range cfgs {
		if cfg.UniqueId == method {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("config is not exist")
}
