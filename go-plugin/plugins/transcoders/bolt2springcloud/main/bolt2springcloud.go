package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/extensions/go-plugin/pkg/protocol/bolt"
	"mosn.io/pkg/protocol/http"
)

const (
	HttpServiceName = "X-TARGET-APP"
	BoltMethodName  = "sofa_head_method_name"
	MosnPath        = "x-mosn-path"
	MosnMethod      = "x-mosn-method"
	MosnHost        = "x-mosn-host"
)

var (
	Http2BoltCode = map[int]uint16{
		http.OK:                  bolt.ResponseStatusSuccess,
		http.BadRequest:          bolt.ResponseStatusError,
		http.InternalServerError: bolt.ResponseStatusServerException,
		http.TooManyRequests:     bolt.ResponseStatusServerThreadpoolBusy,
		http.NotImplemented:      bolt.ResponseStatusNoProcessor,
		http.RequestTimeout:      bolt.ResponseStatusTimeout,
	}
)

type bolt2springcloud struct {
	cfg         map[string]interface{}
	config      *Config
	boltRequest *bolt.Request
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	return &bolt2springcloud{
		cfg: cfg,
	}
}

func (t *bolt2springcloud) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	_, ok := headers.(*bolt.Request)
	if !ok {
		return false
	}
	config, err := t.getConfig(ctx, headers)
	if err != nil {
		return false
	}
	t.config = config
	return true
}

func (t *bolt2springcloud) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceRequest, ok := headers.(*bolt.Request)
	if !ok {
		return headers, buf, trailers, nil
	}
	t.boltRequest = sourceRequest
	// update headers
	targetRequest := t.httpReq2BoltReq(sourceRequest)
	return http.RequestHeader{RequestHeader: targetRequest}, buf, trailers, nil
}

func (t *bolt2springcloud) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceRequest, ok := headers.(http.ResponseHeader)
	if !ok {
		if _, ok := headers.(http.RequestHeader); ok {
			return t.boltRequest, buf, trailers, nil
		}
		return headers, buf, trailers, nil
	}
	targetResponse := bolt.NewRpcResponse(t.boltRequest.RequestId, t.getCode(sourceRequest.StatusCode()), nil, buf)
	targetResponse.Class = t.config.Class
	targetResponse.Codec = bolt.JsonSerialize //json
	return targetResponse, buf, trailers, nil
}

func (t *bolt2springcloud) httpReq2BoltReq(headers *bolt.Request) *fasthttp.RequestHeader {
	targetRequest := &fasthttp.RequestHeader{}
	targetRequest.Set(MosnMethod, t.config.Method)
	targetRequest.Set(MosnPath, t.config.Path)
	targetRequest.Set(HttpServiceName, t.config.TragetApp)
	return targetRequest
}

func (t *bolt2springcloud) getCode(code int) uint16 {
	boltCode, ok := Http2BoltCode[code]
	if ok {
		return boltCode
	}
	return bolt.ResponseStatusUnknown
}

func (t *bolt2springcloud) getConfig(ctx context.Context, headers api.HeaderMap) (*Config, error) {
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
	method, ok := headers.Get(BoltMethodName)
	for _, cfg := range cfgs {
		if cfg.UniqueId == method {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("config is not exist")
}
