package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo"
	"github.com/valyala/fasthttp"
	"math/rand"
	"mosn.io/api"
	"mosn.io/api/extensions/transcoder"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
	"regexp"
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

type paramAdapter struct {
	Service        string           `json:"service"`
	Method         string           `json:"method"`
	Version        string           `json:"version"`
	Group          string           `json:"group"`
	Double         string           `json:"double"`
	Query          []*query         `json:"query"`
	Body           *body            `json:"body"`
	HttpPathParams []*httpPathParam `json:"http_path_params"`
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

var conf = map[string]*paramAdapter{}

//accept return when head has transcoder key and value is equal to TRANSCODER_NAME
func (t *http2dubbo) Accept(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) bool {
	return true
}

// transcode dubbp request to http request
func (t *http2dubbo) TranscodingRequest(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) (api.HeaderMap, api.IoBuffer, api.HeaderMap, error) {

	log.DefaultContextLogger.Debugf(ctx, "[http2dubbo transcoder] request header %v ,buf %v,", headers, buf)
	if httpRequest, ok := headers.(http.RequestHeader); ok {
		service, _ := httpRequest.Get("service")
		method := string(httpRequest.Method())
		path := string(httpRequest.RequestURI())
		items := strings.Split(path, "?")
		param, pathParams := getParamAdapter(catStr(service, ".", items[0], ".", method), conf)
		queryMap := map[string]string{}
		if len(items) == 2 {
			queryMap = getQuery(items[1])
		}
		// 2. assemble target request
		targetRequest, err := EncodeHttp2Dubbo(ctx, headers, buildDubboHttpRequestParams(headers, buf, param, pathParams, queryMap))
		if err != nil {
			return nil, nil, nil, err
		}
		return targetRequest.GetHeader(), targetRequest.GetData(), trailers, nil
	}
	return nil, nil, nil, fmt.Errorf("[http2dubbo transcoder] error for transcode header is not http.RequestHeader")
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
	if err != nil {
		return err
	}
	targetResponse.SetBody(value)
	return nil
}

//encode http require to dubbo require
func EncodeHttp2Dubbo(ctx context.Context, headers api.HeaderMap, reqBody DubboHttpRequestParams) (*dubbo.Frame, error) {
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
	payLoadByteFin, err := EncodeWorkLoad(headers, reqBody)
	if err != nil {
		log.DefaultContextLogger.Errorf(ctx, "[http2dubbo transcoder] error EncodeWorkLoad error %v", err)
		return nil, err
	}
	frame.SetData(buffer.NewIoBufferBytes(payLoadByteFin))

	return frame, nil
}

func EncodeWorkLoad(headers api.HeaderMap, reqBody DubboHttpRequestParams) ([]byte, error) {
	//body
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
	if cfgJson, err := json.Marshal(cfg); err == nil {
		json.Unmarshal(cfgJson, &conf)
	}
	return &http2dubbo{cfg: cfg}
}

func getQuery(queryStr string) map[string]string {
	queryMap := map[string]string{}
	if queryStr == "" {
		return queryMap
	}
	for _, q := range strings.Split(queryStr, "&") {
		m := strings.Split(q, "=")
		if len(m) == 2 {
			queryMap[m[0]] = m[1]
		}
	}

	return queryMap

}

func buildDubboHttpRequestParams(headers api.HeaderMap, buf api.IoBuffer, param *paramAdapter, pathParams []string, queryMap map[string]string) DubboHttpRequestParams {
	reqBody := DubboHttpRequestParams{}
	if param != nil {
		headers.Set("service", param.Service)
		headers.Set("dubbo", param.Double)
		headers.Set("method", param.Method)
		headers.Set("version", param.Version)
		headers.Set("group", param.Group)
		params := []dubbo.Parameter{}

		//路径参数转换
		for i, h := range param.HttpPathParams {
			if i < len(pathParams) {
				v := pathParams[i]
				if v != "" {
					p := dubbo.Parameter{
						Type:  h.Type,
						Value: v,
					}
					params = append(params, p)
				}
			}
		}

		//查询参数转换
		for _, q := range param.Query {
			v := queryMap[q.Key]
			if v != "" {
				p := dubbo.Parameter{
					Type:  q.Type,
					Value: v,
				}
				params = append(params, p)
			}
		}

		//body参数转换
		if param.Body != nil {
			var by interface{}
			json.Unmarshal(buf.Bytes(), &by)
			p := dubbo.Parameter{
				Type:  param.Body.Type,
				Value: by,
			}
			params = append(params, p)
		}
		reqBody.Parameters = params
	}

	return reqBody
}

func catStr(params ...string) string {
	var buffer bytes.Buffer
	for _, v := range params {
		buffer.WriteString(v)
	}
	return buffer.String()
}

func getParamAdapter(uri string, maps map[string]*paramAdapter) (*paramAdapter, []string) {

	param := maps[uri]
	if param != nil {
		return param, nil
	}

	for k, v := range maps {
		flysnowRegexp := regexp.MustCompile(catStr("^", k, "$"))
		params := flysnowRegexp.FindStringSubmatch(uri)

		if params != nil {
			if len(params) > 1 {
				return v, params[1:]
			} else {
				return v, nil
			}
		}
	}

	return nil, nil
}
