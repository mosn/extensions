package bumsbeis

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

type Beis2Bums struct {
	config *Beis2BumsConfig
	ctx    context.Context
	root   *etree.Element
	header api.HeaderMap
}

func NewBeis2Bums(ctx context.Context, header api.HeaderMap, buf api.IoBuffer, config *Beis2BumsConfig) (*Beis2Bums, error) {
	doc := etree.NewDocument()
	doc.ReadFrom(buf)
	root := doc.SelectElement("Document")
	if root == nil {
		return nil, fmt.Errorf("doc is empty:doc :%s", doc.FullTag())
	}
	return &Beis2Bums{
		root:   root,
		header: header,
		config: config,
	}, nil
}

func (bibm *Beis2Bums) BodyJson(header api.HeaderMap) ([]byte, error) {
	val := fastjson.Arena{}
	body := val.NewObject()
	val = fastjson.Arena{}
	head := val.NewObject()
	for _, t := range bibm.root.Child {
		if c, ok := t.(*etree.Element); ok {
			bibm.bodyTrancoder(c, head, body, header)
		}
	}

	val = fastjson.Arena{}
	resp := val.NewObject()
	resp.Set("head", head)
	resp.Set("body", body)
	data := resp.MarshalTo(nil)
	return data, nil
}

func (bibm *Beis2Bums) HeadRespone() (api.HeaderMap, error) {
	respHeader := &fasthttp.ResponseHeader{}
	respHeader.Set("Content-Type", "application/json")
	bibm.header.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			respHeader.Set(key, value)
		}
		return true
	})

	if code, ok := bibm.header.Get("x-mosn-status"); ok {
		statusCode, err := strconv.Atoi(code)
		if err == nil {
			respHeader.SetStatusCode(statusCode)
		} else {
			log.DefaultContextLogger.Warnf(bibm.ctx, "the atoi of statuscode failed. err:%s", err)
		}
	}

	// beis数据解析
	br := bibm.header.(*beis.Response)
	respHeader.Set("VersionId", br.VersionID)
	respHeader.Set("OrigSender", br.OrigSender)
	respHeader.Set("CtrlBits", br.CtrlBits)
	respHeader.Set("AreaCode", br.AreaCode)
	return http.ResponseHeader{respHeader}, nil
}

func (bibm *Beis2Bums) HeadRequest() (api.HeaderMap, error) {
	reqHeader := &fasthttp.RequestHeader{}
	reqHeader.Set("Content-Type", "application/json")
	bibm.header.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			reqHeader.Set(key, value)
		}
		return true
	})
	reqHeader.Set("x-mosn-method", bibm.config.Method)
	reqHeader.Set("x-mosn-path", bibm.config.Path)
	reqHeader.Set("X-TARGET-APP", bibm.config.GWName)

	// beis数据解析
	br := bibm.header.(*beis.Request)
	reqHeader.Set("VersionId", br.VersionID)
	reqHeader.Set("OrigSender", br.OrigSender)
	reqHeader.Set("CtrlBits", br.CtrlBits)
	reqHeader.Set("AreaCode", br.AreaCode)
	return http.RequestHeader{reqHeader}, nil
}

func (bibm *Beis2Bums) Transcoder(isRequest bool) (header api.HeaderMap, buf api.IoBuffer, err error) {
	if isRequest {
		header, err = bibm.HeadRequest()
	} else {
		header, err = bibm.HeadRespone()
	}
	if err != nil {
		return nil, nil, err
	}

	body, err := bibm.BodyJson(header)
	if err != nil {
		return nil, nil, err
	}
	return header, buffer.NewIoBufferBytes(body), nil
}

func (bibm *Beis2Bums) bodyTrancoder(e *etree.Element, head, body *fastjson.Value, header api.HeaderMap) {
	switch e.Tag {
	case "AppHead":
		bibm.BodyHead(e, head, "appHead", header)
	case "SysHead":
		bibm.BodyHead(e, head, "sysHead", header)
	case "LOCAL_HEAD":
		bibm.BodyHead(e, head, "localHead", header)
	case "details":
		bibm.Details(e, body)
	default:
		bibm.OtherBody(e, body, e.Tag, 0)
	}
}

func (bibm *Beis2Bums) BodyHead(e *etree.Element, head *fastjson.Value, key string, header api.HeaderMap) {
	elemEmpty := true
	for _, t := range e.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = true
			bibm.BodyHead(c, head, c.Tag, header)
		}
	}
	if elemEmpty {
		if len(e.Text()) == 0 {
			log.DefaultContextLogger.Warnf(bibm.ctx, "the tag:%s of value is empty in head", e.Tag)
			return
		}
		val := fastjson.Arena{}
		key = bibm.HeadKey(key)
		head.Set(key, val.NewString(e.Text()))
		// TODO 状态码转换映射
		if header != nil {
			// traceid/spanid 兼容
			if strings.EqualFold(key, "traceid") {
				header.Add("SpanId", e.Text())
			}
			if strings.EqualFold(key, "traceid") {
				header.Add("TraceId", e.Text())
			}
		}
	}
}

func (bibm *Beis2Bums) HeadKey(key string) string {
	key = ToFristLower(key)
	switch key {
	case "traceid":
		return "traceId"
	case "spanid":
		return "spanId"
	case "uniqueid":
		return "uniqueId"
	}
	return key
}

func (bibm *Beis2Bums) BodyKey(key string) string {
	return key
	// return ToFristLower(key)
}

func (bibm *Beis2Bums) OtherBody(e *etree.Element, body *fastjson.Value, key string, depth int) {
	elemEmpty := true
	val := fastjson.Arena{}
	obj := val.NewObject()

	for _, t := range e.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = false
			if depth <= 0 {
				bibm.OtherBody(c, obj, c.Tag, depth+1)
			} else {
				bibm.BodyHead(c, obj, c.Tag, nil)
			}
		}
	}

	if elemEmpty {
		if len(e.Text()) == 0 {
			return
		}
		if depth == 1 {
			key = ToFristLower(key)
		}
		body.Set(key, val.NewString(e.Text()))
	} else {
		// 数据拍平，相同key 有覆盖风险,数组只保留最后一组
		body.Set(key, obj)
	}
}

func (bibm *Beis2Bums) Details(e *etree.Element, body *fastjson.Value) {
	val := fastjson.Arena{}
	details := val.NewObject()
	elemEmpty := true
	for _, t := range e.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = false
			bodyKey := bibm.BodyKey(c.Tag)
			array := details.Get(bodyKey)
			if array == nil {
				erena := fastjson.Arena{}
				array = erena.NewArray()
				details.Set(bodyKey, array)
			}
			elements := e.SelectElements(c.Tag)
			bibm.VisitArray(e, array, elements)
			details.Set(bodyKey, array)
		}
	}
	if !elemEmpty {
		body.Set("details", details)
	} else if len(e.Text()) != 0 {
		body.Set("details", val.NewString(e.Text()))
	} else {
		log.DefaultContextLogger.Warnf(bibm.ctx, "the tag:%s of value is empty in details", e.Tag)
	}
}

// 保证数组顺序
func (bibm *Beis2Bums) VisitArray(parent *etree.Element, array *fastjson.Value, elements []*etree.Element) {
	for index, e := range elements {
		array.SetArrayItem(index, bibm.VisitObject(e))
		parent.RemoveChild(e)
	}
}

func (bibm *Beis2Bums) VisitObject(ele *etree.Element) *fastjson.Value {
	erena := fastjson.Arena{}
	object := erena.NewObject()
	elemEmpty := true
	for _, t := range ele.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = false
			if elements := ele.SelectElements(c.Tag); len(elements) == 1 {
				key := bibm.BodyKey(c.Tag)
				object.Set(key, bibm.VisitObject(c))
			} else {
				erena := fastjson.Arena{}
				array := erena.NewArray()
				bibm.VisitArray(ele, array, elements)
				key := bibm.BodyKey(c.Tag)
				object.Set(key, array)
			}
		}
	}
	if elemEmpty {
		val := fastjson.Arena{}
		if len(ele.Text()) == 0 {
			log.DefaultContextLogger.Warnf(bibm.ctx, "the tag:%s of value is empty in VisitObject", ele.Tag)
		}
		return val.NewString(ele.Text())
	}
	return object
}
