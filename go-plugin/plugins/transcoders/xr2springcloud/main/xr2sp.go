/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	"mosn.io/api/extensions/transcoder"
	"mosn.io/extensions/go-plugin/pkg/protocol/xr"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/protocol/http"
)

type xr2springcloud struct {
	cfg    map[string]interface{}
	config *Config
}

func (t *xr2springcloud) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	_, ok := headers.(*xr.Request)
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

func (t *xr2springcloud) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceHeader, ok := headers.(*xr.Request)
	if !ok {
		return nil, nil, nil, fmt.Errorf("[xprotocol][xr] decode xr header type error")
	}
	reqHeaderImpl := &fasthttp.RequestHeader{}
	sourceHeader.Header.Range(func(key, value string) bool {
		if key != fasthttp.HeaderContentLength {
			reqHeaderImpl.Set(key, value)
		}
		return true
	})
	reqHeaderImpl.Set("x-mosn-method", t.config.Method)
	reqHeaderImpl.Set("x-mosn-path", t.config.Path)
	reqHeaderImpl.Set("X-TARGET-APP", t.config.TragetApp)
	reqHeaders := http.RequestHeader{reqHeaderImpl}
	return reqHeaders, buf, trailers, nil
}

func (t *xr2springcloud) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceHeader, ok := headers.(http.ResponseHeader)
	if !ok {
		return nil, nil, nil, fmt.Errorf("[xprotocol][xr] decode http header type error")
	}
	//header
	xrResponse := xr.Response{}
	sourceHeader.Range(func(key, value string) bool {
		//skip for Content-Length,the Content-Length may effect the value decode when transcode more one time
		if key != "Content-Length" && key != "Accept:" {
			xrResponse.Set(key, value)
		}
		return true
	})

	payloads := buffer.NewIoBufferBytes(buf.Bytes())
	respHeader := xr.NewRpcResponse(&xrResponse.Header, payloads)
	if respHeader == nil {
		return nil, nil, nil, fmt.Errorf("[xprotocol][xr] decode http header type error")
	}
	return respHeader.GetHeader(), respHeader.GetData(), trailers, nil
}

func LoadTranscoderFactory(cfg map[string]interface{}) transcoder.Transcoder {
	return &xr2springcloud{cfg: cfg}
}
