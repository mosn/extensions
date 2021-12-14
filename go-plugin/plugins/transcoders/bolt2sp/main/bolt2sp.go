package main

import (
	"context"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/bolt"
	"github.com/valyala/fasthttp"
	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/pkg/protocol/http"
)

const (
	DefalutClass = "com.alipay.sofa.rpc.core.response.SofaResponse"
	ClassName    = "class"
)

var (
	HeaderKeys    = []string{"x-mosn-host", "x-mosn-method", "x-mosn-path"}
	Http2BoltCode = map[int]uint16{
		http.OK:                  bolt.ResponseStatusSuccess,
		http.BadRequest:          bolt.ResponseStatusError,
		http.InternalServerError: bolt.ResponseStatusServerException,
		http.TooManyRequests:     bolt.ResponseStatusServerThreadpoolBusy,
		http.NotImplemented:      bolt.ResponseStatusNoProcessor,
		http.RequestTimeout:      bolt.ResponseStatusTimeout,
	}
)

type bolt2sp struct {
	cfg         map[string]interface{}
	boltRequest *bolt.Request
}

func LoadTranscoderFactory(cfg map[string]interface{}) at.Transcoder {
	return &bolt2sp{
		cfg: cfg,
	}
}

func (t *bolt2sp) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

func (t *bolt2sp) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceRequest, ok := headers.(*bolt.Request)
	if !ok {
		return headers, buf, trailers, nil
	}
	t.boltRequest = sourceRequest
	targetRequest := fasthttp.Request{}
	sourceRequest.Range(func(Key, Value string) bool {
		targetRequest.Header.Set(Key, Value)
		return true
	})
	// update headers
	targetRequest = t.updateRequestHeader(targetRequest)
	return http.RequestHeader{RequestHeader: &targetRequest.Header}, buf, trailers, nil
}

func (t *bolt2sp) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceRequest, ok := headers.(http.ResponseHeader)
	if !ok {
		return headers, buf, trailers, nil
	}
	targetResponse := bolt.NewRpcResponse(t.boltRequest.RequestId, t.getCode(sourceRequest.StatusCode()), nil, buf)
	targetResponse.Class = t.getClass()
	targetResponse.Codec = bolt.JsonSerialize //json
	return targetResponse, buf, trailers, nil
}

func (t *bolt2sp) updateRequestHeader(req fasthttp.Request) fasthttp.Request {
	if t.cfg == nil {
		return req
	}
	for _, key := range HeaderKeys {
		val, ok := t.cfg[key].(string)
		if ok && val != "" {
			req.Header.Set(key, val)
		}
	}
	return req
}

func (t *bolt2sp) getClass() string {
	if t.cfg == nil {
		return DefalutClass
	}
	name, ok := t.cfg[ClassName]
	if ok {
		return DefalutClass
	}
	n, ok := name.(string)
	if ok {
		return DefalutClass
	}
	return n
}

func (t *bolt2sp) getCode(code int) uint16 {
	boltCode, ok := Http2BoltCode[code]
	if ok {
		return boltCode
	}
	return bolt.ResponseStatusUnknown
}
