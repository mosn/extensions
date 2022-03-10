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
	"mosn.io/extensions/go-plugin/pkg/protocol/cd"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/protocol/http"
)

type xr2sp struct{ cfg map[string]interface{} }

func (t *xr2sp) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

func (t *xr2sp) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceHeader, ok := headers.(*cd.Request)
	if !ok {
		return nil, nil, nil, fmt.Errorf("[xprotocol][dubbo] decode xr header type error")
	}
	reqHeaderImpl := &fasthttp.RequestHeader{}
	sourceHeader.Header.Range(func(key, value string) bool {
		if key != fasthttp.HeaderContentLength {
			reqHeaderImpl.Set(key, value)
		}
		return true
	})
	//set request idoujju
	serviceCode, _ := headers.Get("ServiceCode")
	reqHeaderImpl.Set("service", serviceCode)
	reqHeaders := http.RequestHeader{reqHeaderImpl}

	return reqHeaders, buf, trailers, nil
}

func (t *xr2sp) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	sourceHeader, ok := headers.(http.ResponseHeader)
	if !ok {
		return nil, nil, nil, fmt.Errorf("[xprotocol][dubbo] decode http header type error")
	}
	//header
	xrResponse := cd.Response{}
	sourceHeader.Range(func(key, value string) bool {
		//skip for Content-Length,the Content-Length may effect the value decode when transcode more one time
		if key != "Content-Length" && key != "Accept:" {
			xrResponse.Set(key, value)
		}
		return true
	})

	payloads := buffer.NewIoBufferBytes(buf.Bytes())
	respHeader := cd.NewRpcResponse(&xrResponse.Header, payloads)

	return respHeader.GetHeader(), respHeader.GetData(), trailers, nil
}

func LoadTranscoderFactory(cfg map[string]interface{}) transcoder.Transcoder {
	return &xr2sp{cfg: cfg}
}
