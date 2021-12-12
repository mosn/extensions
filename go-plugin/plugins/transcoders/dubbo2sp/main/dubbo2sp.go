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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mosn.io/api"
	"mosn.io/api/extensions/transcoder"
	"strconv"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo"
	"github.com/valyala/fasthttp"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/protocol/http"
)

const HTTP_DUBBO_REQUEST_ID_NAME = "Dubbo-Request-Id"

type DubboHttpResponseBody struct {
	Attachments map[string]string `json:"attachments"`
	Value       interface{}       `json:"value"`
	Exception   string            `json:"exception"`
}

type dubbo2http struct{ cfg map[string]interface{} }

//accept return when head has transcoder key and value is equal to TRANSCODER_NAME
func (t *dubbo2http) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

//dubbo request 2 http request
func (t *dubbo2http) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	// 1. set sub protocol
	sourceHeader, ok := headers.(*dubbo.Frame)
	if !ok {
		return nil, nil, nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo header type error")
	}
	//// 2. assemble target request
	byteData, err := DeocdeWorkLoad(headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	reqHeaderImpl := &fasthttp.RequestHeader{}
	sourceHeader.Header.CommonHeader.Range(func(key, value string) bool {
		if key != fasthttp.HeaderContentLength {
			reqHeaderImpl.SetCanonical([]byte(key), []byte(value))
		}
		return true
	})
	//set request id
	reqHeaderImpl.Set(HTTP_DUBBO_REQUEST_ID_NAME, strconv.FormatUint(sourceHeader.Id, 10))
	reqHeaders := http.RequestHeader{reqHeaderImpl}
	return reqHeaders, buffer.NewIoBufferBytes(byteData), nil, nil
}

// encode the dubbo request body 2 http request body
func DeocdeWorkLoad(headers api.HeaderMap, buf api.IoBuffer) ([]byte, error) {
	var paramsTypes string
	sourceRequest, ok := headers.(*dubbo.Frame)
	if !ok {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode header type error")
	}
	dataArr := bytes.Split(sourceRequest.GetData().Bytes(), []byte{10})
	err := json.Unmarshal(dataArr[4], &paramsTypes)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode params fail")
	}
	count := dubbo.GetArgumentCount(paramsTypes)
	//skip useless dubbo flags
	arrs := dataArr[5 : 5+count]
	// decode dubbo body
	params, err := dubbo.DecodeParams(paramsTypes, arrs)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode params fail")
	}
	attachments := map[string]string{}
	err = json.Unmarshal(dataArr[5+count], &attachments)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode attachments fail")
	}
	//encode to http budy
	content := map[string]interface{}{}
	content["attachments"] = attachments
	content["parameters"] = params
	byte, _ := json.Marshal(content)
	return byte, nil
}

//http2dubbo
func (t *dubbo2http) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	targetRequest, err := DecodeHttp2Dubbo(headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	return targetRequest.GetHeader(), targetRequest.GetData(), trailers, nil
}

// decode http response to dubbo response
func DecodeHttp2Dubbo(headers api.HeaderMap, buf api.IoBuffer) (*dubbo.Frame, error) {

	sourceHeader, ok := headers.(http.ResponseHeader)
	if !ok {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo header type error")
	}
	//header
	allHeaders := map[string]string{}
	sourceHeader.Range(func(key, value string) bool {
		//skip for Content-Length,the Content-Length may effect the value decode when transcode more one time
		if key != "Content-Length" && key != "Accept:" {
			allHeaders[key] = value
		}
		return true
	})
	frame := &dubbo.Frame{
		Header: dubbo.Header{
			CommonHeader: allHeaders,
		},
	}
	// convert data to dubbo frame
	workLoad, err := EncodeWorkLoad(headers, buf)
	if err != nil {
		return nil, err
	}

	//magic
	frame.Magic = dubbo.MagicTag

	//  flag
	frame.Flag = 0x46
	// status when http return not ok, return error
	if sourceHeader.StatusCode() != http.OK {
		//BAD_RESPONSE
		frame.Status = 40
	} else {
		frame.Status = 20
	}
	// decode request id
	if id, ok := allHeaders[HTTP_DUBBO_REQUEST_ID_NAME]; !ok {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo id missed")
	} else {
		frameId, _ := strconv.ParseInt(id, 10, 64)
		frame.Id = uint64(frameId)
	}

	// event
	frame.IsEvent = false
	// twoway
	frame.IsTwoWay = true
	// direction
	frame.Direction = dubbo.EventResponse
	// serializationId json
	frame.SerializationId = 6
	frame.SetData(buffer.NewIoBufferBytes(workLoad))
	return frame, nil
}

// http response body example: {"attachments":null,"value":"22222","exception":null}
// make dubbo workload
func EncodeWorkLoad(headers api.HeaderMap, buf api.IoBuffer) ([]byte, error) {
	responseBody := DubboHttpResponseBody{}
	workload := [][]byte{}
	if buf == nil {
		return nil, fmt.Errorf("no workload to decode")
	}
	err := json.Unmarshal(buf.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}
	if responseBody.Exception == "" {
		//out.writeByte(
		if responseBody.Value == nil {
			resType, _ := json.Marshal(dubbo.RESPONSE_NULL_VALUE)
			workload = append(workload, resType)
		} else {
			resType, _ := json.Marshal(dubbo.RESPONSE_VALUE)
			workload = append(workload, resType)
			ret, _ := json.Marshal(responseBody.Value)
			workload = append(workload, ret)

		}
	} else {
		resType, _ := json.Marshal(dubbo.RESPONSE_WITH_EXCEPTION)
		workload = append(workload, resType)
		ret, _ := json.Marshal(responseBody.Exception)
		workload = append(workload, ret)
	}
	workloadByte := bytes.Join(workload, []byte{10})

	return workloadByte, nil
}

func LoadTranscoderFactory(cfg map[string]interface{}) transcoder.Transcoder {
	return &dubbo2http{cfg: cfg}
}
