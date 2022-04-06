package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

type Config struct {
	UniqueId    string      `json:"unique_id"`
	TragetApp   string      `json:"target_app"`
	ReqMapping  ReqMapping  `json:"req_mapping"`
	RespMapping RespMapping `json:"resp_mapping"`
}

type ReqMapping struct {
	Version    string           `json:"version"`
	Group      string           `json:"group"`
	Double     string           `json:"double"`
	Method     string           `json:"method"`
	Query      []*query         `json:"query"`
	Body       *body            `json:"body"`
	PathParams []*httpPathParam `json:"path_params"`
	pathParams []string         `json:"-"`
}

type RespMapping struct {
}

type query struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type httpPathParam struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type body struct {
	Type string `json:"type"`
}

func (t *springcloud2dubbo) getConfig(ctx context.Context, headers *fasthttp.RequestHeader) (*Config, error) {
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
	service := string(headers.RequestURI())
	index := strings.Index(service, "?")
	if index != -1 {
		service = service[:index]
	}
	method := string(headers.Method())
	method = service + "." + method
	for _, cfg := range cfgs {
		if cfg.UniqueId == method {
			return cfg, nil
		}
	}

	for _, cfg := range cfgs {
		flysnowRegexp := regexp.MustCompile(catStr("^", cfg.UniqueId, "$"))
		params := flysnowRegexp.FindStringSubmatch(method)
		if params != nil {
			if len(params) > 1 {
				cfg.ReqMapping.pathParams = params[1:]
				return cfg, nil
			} else {
				return cfg, nil
			}
		}
	}
	return nil, fmt.Errorf("config is not exist")
}
