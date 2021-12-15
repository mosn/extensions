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
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo"
	"github.com/valyala/fasthttp"
	"mosn.io/api"
	"mosn.io/api/extensions/transcoder"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
	"strconv"
)

const HTTP_DUBBO_REQUEST_ID_NAME = "Dubbo-Request-Id"

type DubboHttpResponseBody struct {
	Attachments map[string]string `json:"attachments"`
	Value       interface{}       `json:"value"`
	Exception   string            `json:"exception"`
}

//"service": "com.alipay.sofa.ms.service.EchoService",
//"http_path": "/reservations/echo",
//"http_method": "POST",
//"http_service": "reservation-service"
type paramAdapter struct {
	HttpPath    string   `json:"http_path"`
	HttpMethod  string   `json:"http_method"`
	HttpService string   `json:"http_service"`
	HttpQuery   []*query `json:"http_query"`
	HttpBody    *body    `json:"http_body"`
}

type query struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type body struct {
	Type string `json:"type"`
}

var conf = map[string]*paramAdapter{}

type dubbo2http struct {
	cfg map[string]interface{}
	Id  uint64
}

//accept return when head has transcoder key and value is equal to TRANSCODER_NAME
func (t *dubbo2http) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

//dubbo request 2 http request
func (t *dubbo2http) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {

	log.DefaultContextLogger.Debugf(ctx, "[dubbo2http transcoder] request header %v ,buf %v,", headers, buf)
	// 1. set sub protocol
	sourceHeader, ok := headers.(*dubbo.Frame)
	if !ok {
		return nil, nil, nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo header type error")
	}
	t.Id = sourceHeader.GetRequestId()
	//// 2. assemble target request
	content, err := DeocdeWorkLoad(headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	reqHeaderImpl := &fasthttp.RequestHeader{}
	sourceHeader.Header.CommonHeader.Range(func(key, value string) bool {
		if key != fasthttp.HeaderContentLength {
			reqHeaderImpl.Set(key, value)
		}
		return true
	})
	service, _ := sourceHeader.Get("service")
	method, _ := sourceHeader.Get("method")
	param := conf[service+"."+method]

	querys := map[string]string{}
	var byteData []byte
	if param != nil {
		reqHeaderImpl.Set("x-mosn-method", param.HttpMethod)
		reqHeaderImpl.Set("service", param.HttpService)
		if params, ok := content["parameters"].([]dubbo.Parameter); ok {
			for _, p := range params {
				if param.HttpBody != nil && p.Type == param.HttpBody.Type {
					byteData, _ = json.Marshal(p.Value)
				}
				for _, q := range param.HttpQuery {
					if p.Type == q.Type {
						if querys[q.Key] == "" {
							querys[q.Key] = Strval(p.Value)
							break
						}
					}
				}
			}

		}

		queryStr := ""
		for k, v := range querys {
			queryStr = queryStr + "&" + k + "=" + v
		}

		path := param.HttpPath
		if queryStr != "" {
			reqHeaderImpl.Set("x-mosn-querystring", queryStr[1:])
		}
		reqHeaderImpl.Set("x-mosn-path", path)
	}

	//set request id
	//reqHeaderImpl.Set(HTTP_DUBBO_REQUEST_ID_NAME, strconv.FormatUint(sourceHeader.Id, 10))
	reqHeaderImpl.Set("Content-Type", "application/json")
	reqHeaders := http.RequestHeader{reqHeaderImpl}
	return reqHeaders, buffer.NewIoBufferBytes(byteData), nil, nil
}

// encode the dubbo request body 2 http request body
func DeocdeWorkLoad(headers api.HeaderMap, buf api.IoBuffer) (map[string]interface{}, error) {
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
	return content, nil
}

//http2dubbo
func (t *dubbo2http) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	log.DefaultContextLogger.Debugf(ctx, "[dubbo2http transcoder] response header %v ,buf %v,", headers, buf)
	targetRequest, err := DecodeHttp2Dubbo(headers, buf, t.Id)
	if err != nil {
		return nil, nil, nil, err
	}
	return targetRequest.GetHeader(), targetRequest.GetData(), trailers, nil
}

// decode http response to dubbo response
func DecodeHttp2Dubbo(headers api.HeaderMap, buf api.IoBuffer, id uint64) (*dubbo.Frame, error) {

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
	workLoad, err := EncodeWorkLoad(sourceHeader, buf)
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
	frame.Id = id

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
func EncodeWorkLoad(headers http.ResponseHeader, buf api.IoBuffer) ([]byte, error) {
	responseBody := DubboHttpResponseBody{}
	workload := [][]byte{}
	if buf == nil {
		return nil, fmt.Errorf("no workload to decode")
	}

	if headers.StatusCode() >= 400 {
		resType, _ := json.Marshal(dubbo.RESPONSE_WITH_EXCEPTION)
		workload = append(workload, resType)
		ret, _ := json.Marshal(responseBody.Exception)
		workload = append(workload, ret)
	} else {
		if buf == nil {
			resType, _ := json.Marshal(dubbo.RESPONSE_NULL_VALUE)
			workload = append(workload, resType)
		} else {
			resType, _ := json.Marshal(dubbo.RESPONSE_VALUE)
			workload = append(workload, resType)
			workload = append(workload, buf.Bytes())

		}
	}

	//err := json.Unmarshal(buf.Bytes(), &responseBody)
	//if err != nil {
	//	return nil, err
	//}
	//if responseBody.Exception == "" {
	//	//out.writeByte(
	//	if responseBody.Value == nil {
	//		resType, _ := json.Marshal(dubbo.RESPONSE_NULL_VALUE)
	//		workload = append(workload, resType)
	//	} else {
	//		resType, _ := json.Marshal(dubbo.RESPONSE_VALUE)
	//		workload = append(workload, resType)
	//		ret, _ := json.Marshal(responseBody.Value)
	//		workload = append(workload, ret)
	//
	//	}
	//} else {
	//	resType, _ := json.Marshal(dubbo.RESPONSE_WITH_EXCEPTION)
	//	workload = append(workload, resType)
	//	ret, _ := json.Marshal(responseBody.Exception)
	//	workload = append(workload, ret)
	//}

	workloadByte := bytes.Join(workload, []byte{10})

	return workloadByte, nil
}

func LoadTranscoderFactory(cfg map[string]interface{}) transcoder.Transcoder {

	if cfgJson, err := json.Marshal(cfg); err == nil {
		json.Unmarshal(cfgJson, &conf)
	}

	return &dubbo2http{cfg: cfg}
}

func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
