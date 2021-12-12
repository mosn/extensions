package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo"
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo/common"
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo/constants"
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo/govern_util"
	"github.com/valyala/fasthttp"
	"math/rand"
	"mosn.io/api"
	"mosn.io/api/extensions/transcoder"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
	"strconv"
	"strings"
)

type http2dubbo struct{ cfg map[string]interface{} }

type DubboHttpResponseBody struct {
	Attachments map[string]string `json:"attachments"`
	Value       interface{}       `json:"value"`
	Exception   string            `json:"exception"`
}

type DubboHttpRequestParams struct {
	Attachments map[string]string `json:"attachments"`
	Parameters  []dubbo.Parameter `json:"parameters"`
}

//accept return when head has transcoder key and value is equal to TRANSCODER_NAME
func (t *http2dubbo) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

// transcode dubbp request to http request
func (t *http2dubbo) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {

	// 2. assemble target request
	targetRequest, err := EncodeHttp2Dubbo(ctx, headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	return targetRequest.GetHeader(), targetRequest.GetData(), trailers, nil
}

// transcode dubbo response to http response
func (t *http2dubbo) TranscodingResponse(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {
	log.DefaultContextLogger.Debugf(ctx, "[http2dubbo transcoder] response header %v ,buf %v,", headers, buf)
	response, err := DecodeDubbo2Http(ctx, headers, buf, trailers)
	if err != nil {
		return nil, nil, nil, err
	}
	return http.ResponseHeader{ResponseHeader: &response.Header}, buffer.NewIoBufferBytes(response.Body()), trailers, nil
}

// decode dubbo response to http response
func DecodeDubbo2Http(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (fasthttp.Response, error) {
	sourceResponse, ok := headers.(*dubbo.Frame)
	if !ok {
		return fasthttp.Response{}, fmt.Errorf("[xprotocol][http] decode http header type error")
	}
	targetResponse := fasthttp.Response{}
	//head
	err := setResponseHeader(ctx, sourceResponse, &targetResponse)
	if err != nil {
		return targetResponse, err
	}
	//body

	if err := setTargetBody(sourceResponse, &targetResponse); err != nil {
		return targetResponse, err
	}

	return targetResponse, nil
}

func setResponseHeader(ctx context.Context, sourceResponse *dubbo.Frame, targetResponse *fasthttp.Response) error {
	// 1. headers
	sourceResponse.Range(func(key, Value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			targetResponse.Header.Set(key, Value)
		}
		return true
	})
	// is fream response
	if sourceResponse.Direction != dubbo.EventResponse {
		log.DefaultContextLogger.Errorf(ctx, "[http2dubbo transcoder] error for transcode header, sourceResponse: %v is not a response", sourceResponse)
		return fmt.Errorf("[http2dubbo transcoder] error for transcode header, sourceResponse: %v is not a response", sourceResponse)
	}
	if code, ok := sourceResponse.Get("x-mosn-status"); ok {
		log.DefaultContextLogger.Debugf(ctx, "[http2dubbo transcoder] get %v code is %v", "x-mosn-status", code)
		statusCode, err := strconv.Atoi(code)
		if err != nil {
			log.DefaultContextLogger.Errorf(ctx, "[http2dubbo transcoder] error for source response header name: %v code: %v, error %v", "x-mosn-status", code, err)
			return fmt.Errorf("[http2dubbo transcoder] error for source response header name: %v code: %v, error %v", "x-mosn-status", code, err)
		}
		targetResponse.SetStatusCode(statusCode)
	}
	if code, _ := sourceResponse.Header.Get("X-Govern-Resp-Code"); code != "" {
		resCode, err := strconv.Atoi(code)
		if err != nil {
			log.DefaultContextLogger.Errorf(ctx, "[http2dubbo transcoder] error for source response header name: %v code: %v, error %v", "X-Govern-Resp-Code", code, err)
			return fmt.Errorf("[http2dubbo transcoder] error for source response name: %v code: %v, error %v", "X-Govern-Resp-Code", code, err)
		}
		targetResponse.SetStatusCode(resCode)
	}
	return nil
}

func setTargetBody(sourceResponse *dubbo.Frame, targetResponse *fasthttp.Response) error {
	//body
	dataArr := bytes.Split(sourceResponse.GetData().Bytes(), []byte{10})
	var resType int
	err := json.Unmarshal(dataArr[0], &resType)
	if err != nil {
		return err
	}
	exception, value, attachments := []byte(`""`), []byte(`""`), []byte(`{}`)
	switch resType {
	case dubbo.RESPONSE_WITH_EXCEPTION:
		// error
		exception = dataArr[1]
	case dubbo.RESPONSE_VALUE:
		value = dataArr[1]
	case dubbo.RESPONSE_NULL_VALUE:
	case dubbo.RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS:
		exception = dataArr[1]
		value = dataArr[3]
	case dubbo.RESPONSE_VALUE_WITH_ATTACHMENTS:
		value = dataArr[1]
		attachments = dataArr[3]
	case dubbo.RESPONSE_NULL_VALUE_WITH_ATTACHMENTS:
		attachments = dataArr[2]
	}
	responseBody := DubboHttpResponseBody{}

	if value != nil && len(value) > 0 {
		if err := json.Unmarshal(value, &responseBody.Value); err != nil {
			return err
		}
	}
	if exception != nil && len(exception) > 0 {
		if err := json.Unmarshal(exception, &responseBody.Exception); err != nil {
			return err
		}
	}
	if attachments != nil && len(attachments) > 0 {
		if err := json.Unmarshal(attachments, &responseBody.Attachments); err != nil {
			return err
		}
	}
	body, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}
	targetResponse.SetBody(body)
	return nil
}

//encode http require to dubbo require
func EncodeHttp2Dubbo(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer) (*dubbo.Frame, error) {
	//process govern filter
	headers, err := setGovernHead(ctx, headers, buf)
	if err != nil {
		return nil, err

	}
	//header
	allHeaders := map[string]string{}
	headers.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			allHeaders[key] = value
		}
		return true
	})
	// convert data to dubbo frame
	frame := &dubbo.Frame{
		Header: dubbo.Header{
			CommonHeader: allHeaders,
		},
	}
	//magic
	frame.Magic = dubbo.MagicTag
	//  flag
	frame.Flag = 0xc6
	// status
	frame.Status = 0
	// decode request id
	frame.Id = rand.Uint64()
	// event
	frame.IsEvent = (frame.Flag & (1 << 5)) != 0
	// twoway
	frame.IsTwoWay = (frame.Flag & (1 << 6)) != 0
	frame.Direction = dubbo.EventRequest
	// serializationId
	frame.SerializationId = int(frame.Flag & 0x1f)
	//workload
	payLoadByteFin, err := EncodeWorkLoad(headers, buf)
	if err != nil {
		log.DefaultContextLogger.Errorf(ctx, "[http2dubbo transcoder] error EncodeWorkLoad error %v", err)
		return nil, err
	}
	frame.SetData(buffer.NewIoBufferBytes(payLoadByteFin))

	return frame, nil
}

func setGovernHead(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer) (api.HeaderMap, error) {
	//X-Govern-Service
	if serviceName, ok := headers.Get(constants.XPROTOCOL_DUBBO_META_INTERFACE); ok {
		if !strings.Contains(serviceName, "@dubbo") {
			version, _ := headers.Get(constants.XPROTOCOL_DUBBO_META_VERSION)
			if version == "0.0.0" {
				version = ""
			}
			group, _ := headers.Get(constants.XPROTOCOL_DUBBO_META_GROUP)
			id := common.BuildDubboDataId(serviceName, version, group)
			govern_util.AddGovernValue(nil, headers, govern_util.GOVERN_SERVICE_KEY, id)
		}
	} else {
		log.DefaultContextLogger.Debugf(ctx, "[stream filter][transcoder] cant read %s in headers", constants.XPROTOCOL_DUBBO_META_INTERFACE)
		return nil, fmt.Errorf("transcoder error: cant read %s in headers ", constants.XPROTOCOL_DUBBO_META_INTERFACE)
	}
	//X-Govern-Service-Type
	govern_util.AddGovernValue(ctx, headers, govern_util.GOVERN_SERVICE_TYPE_KEY, constants.XPROTOCOL_TYPE_DUBBO)

	//X-Govern-Method
	if method, ok := headers.Get(dubbo.MethodNameHeader); ok {
		govern_util.AddGovernValue(ctx, headers, dubbo.MethodNameHeader, method)
	} else {
		log.DefaultContextLogger.Debugf(ctx, "[stream filter][transcoder] cant read %s in headers", dubbo.MethodNameHeader)
		return nil, fmt.Errorf("transcoder error: cant read %s in headers ", dubbo.MethodNameHeader)
	}
	//X-Govern-Timeout
	if val, ok := headers.Get("timeout"); ok {
		govern_util.AddGovernValue(ctx, headers, govern_util.GOVERN_TIMEOUT_KEY, val)
	}
	//X-Govern-Source-App
	if sourceApp, ok := headers.Get(constants.SOFA_RPC_HEADER_CALLER_APP_KEY); ok {
		govern_util.AddGovernValue(ctx, headers, govern_util.GOVERN_SOURCE_APP_KEY, sourceApp)
	}
	//X-Govern-Target-App
	if targetApp, ok := headers.Get(constants.SOFARPC_ROUTER_HEADER_TARGET_APP_KEY); ok {
		govern_util.AddGovernValue(ctx, headers, constants.SOFARPC_ROUTER_HEADER_TARGET_APP_KEY, targetApp)
	}

	return headers, nil
}

func EncodeWorkLoad(headers api.HeaderMap, buf api.IoBuffer) ([]byte, error) {
	//body
	var reqBody DubboHttpRequestParams
	if buf == nil {
		return nil, fmt.Errorf("nil buf data error")
	}
	if err := json.Unmarshal(buf.Bytes(), &reqBody); err != nil {
		return nil, err
	}
	if reqBody.Attachments == nil {
		reqBody.Attachments = make(map[string]string)
	}

	//设置playload
	headers.Range(func(key, value string) bool {
		reqBody.Attachments[key] = value
		return true
	})

	//service
	serviceName := dubbo.HeadGetDefault(headers, "service", "")
	index := strings.Index(serviceName, "@")
	if index > 0 {
		serviceName = serviceName[:index]
	}
	index = strings.Index(serviceName, ":")
	if index > 0 {
		serviceName = serviceName[:index]
	}
	reqBody.Attachments["interface"] = serviceName

	dubboVersion := dubbo.HeadGetDefault(headers, "dubbo", "2.6.5")
	serviceVersion := dubbo.HeadGetDefault(headers, "version", "0.0.0")
	reqBody.Attachments["version"] = serviceVersion

	serviceMethod := dubbo.HeadGetDefault(headers, "method", "")
	serviceGroup := dubbo.HeadGetDefault(headers, "group", "")
	if serviceGroup != "" {
		reqBody.Attachments["group"] = serviceGroup
	}

	dubboVersionByte := []byte(`"` + dubboVersion + `"`)
	serviceNameByte := []byte(`"` + serviceName + `"`)
	verionByte := []byte(`"` + serviceVersion + `"`)
	methodNameByte := []byte(`"` + serviceMethod + `"`)

	//有几个类型写一个类型，直接跟在字符串后面，数组的前面加[，基本类型要转义,不换行
	//有几个参数写几个参数，要写在byte[]后面,换行
	paramesByte := []byte{}
	var paramesTypeStr string
	if reqBody.Parameters != nil || len(reqBody.Parameters) > 0 {
		for i := 0; i < len(reqBody.Parameters); i++ {
			if reqBody.Parameters[i].Type != "" {
				paramesTypeStr += dubbo.EncodeRequestType(reqBody.Parameters[i].Type)
				valByte, _ := json.Marshal(reqBody.Parameters[i].Value)
				paramesByte = append(paramesByte, valByte...)
				if i < len(reqBody.Parameters)-1 {
					paramesByte = append(paramesByte, []byte{10}...)
				}
			}
		}
	}

	var paramesTypeByte []byte
	var attachmentsByte []byte
	if paramesTypeStr != "" {
		paramesTypeByte, _ = json.Marshal(paramesTypeStr)
		attachmentsByte, _ = json.Marshal(reqBody.Attachments)
	}

	payLoadByteFin := bytes.Join([][]byte{dubboVersionByte, serviceNameByte, verionByte, methodNameByte, paramesTypeByte, paramesByte, attachmentsByte}, []byte{10})

	return payLoadByteFin, nil
}

func LoadTranscoderFactory(cfg map[string]interface{}) transcoder.Transcoder {
	return &http2dubbo{cfg: cfg}
}
