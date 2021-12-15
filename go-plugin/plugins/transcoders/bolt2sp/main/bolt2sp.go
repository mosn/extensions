package main

import (
	"context"
	"errors"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/bolt"
	"github.com/valyala/fasthttp"
	"mosn.io/api"
	at "mosn.io/api/extensions/transcoder"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

const (
	DefaultClass   = "com.alipay.sofa.rpc.core.response.SofaResponse"
	DefaultPath    = "/"
	ClassName      = "class"
	ServiceName    = "service"
	BoltMethodName = "sofa_head_method_name"
	MosnPath       = "x-mosn-path"
	MosnMethod     = "x-mosn-method"
	MosnHost       = "x-mosn-host"
)

var (
	ErrEmptyPath    = errors.New("path is empty")
	ErrEmptyService = errors.New("service is empty")
	HttpHeaderKeys  = []string{"x-mosn-host", "x-mosn-method", "x-mosn-path"}
	Http2BoltCode   = map[int]uint16{
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
	_, ok := headers.(*bolt.Request)
	return ok
}

func (t *bolt2sp) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceRequest, ok := headers.(*bolt.Request)
	if !ok {
		return headers, buf, trailers, nil
	}
	t.boltRequest = sourceRequest
	// update headers
	targetRequest := t.httpReq2BoltReq(sourceRequest)
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

func (t *bolt2sp) httpReq2BoltReq(headers *bolt.Request) fasthttp.Request {
	targetRequest := fasthttp.Request{}
	headers.Range(func(Key, Value string) bool {
		targetRequest.Header.Set(Key, Value)
		return true
	})

	if t.cfg == nil {
		targetRequest.Header.Set(MosnPath, DefaultPath)
		return targetRequest
	}
	if err := t.updateHttpPath(headers, &targetRequest); err != nil {
		log.DefaultLogger.Warnf("[bolt2s] [httpReq2BoltReq] parse conf failed:%v", err)
	}
	if err := t.updateHttpService(&targetRequest); err != nil {
		log.DefaultLogger.Warnf("[bolt2s] [UpdateHttpService] failed:%v", err)
	}
	return targetRequest
}

//service

func (t *bolt2sp) updateHttpService(httpHeaders *fasthttp.Request) error {
	value, ok := t.cfg[ServiceName].(string)
	if !ok {
		return ErrEmptyService
	}
	httpHeaders.Header.Set(ServiceName, value)
	return nil
}

func (t *bolt2sp) updateHttpPath(boltHeaders *bolt.Request, httpHeaders *fasthttp.Request) error {
	service, ok := boltHeaders.Get(ServiceName)
	if !ok {
		return ErrEmptyPath
	}
	method, ok := boltHeaders.Get(BoltMethodName)
	if !ok {
		return ErrEmptyPath
	}
	serviceMap, ok := t.cfg[service].(map[string]interface{})
	if !ok {
		return ErrEmptyPath
	}
	datas, ok := serviceMap[method].(map[string]interface{})
	if !ok {
		return ErrEmptyPath
	}
	for _, key := range HttpHeaderKeys {
		value, ok := datas[key].(string)
		if ok {
			httpHeaders.Header.Set(key, value)
		}
	}
	return nil
}

func (t *bolt2sp) getClass() string {
	if t.cfg == nil {
		return DefaultClass
	}
	name, ok := t.cfg[ClassName]
	if ok {
		return DefaultClass
	}
	n, ok := name.(string)
	if ok {
		return DefaultClass
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
