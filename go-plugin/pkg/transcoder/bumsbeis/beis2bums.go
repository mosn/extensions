package bumsbeis

import (
	"context"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
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
	//TODO xml msg to serviceCode，senceCode
	return header, buffer.NewIoBufferBytes(body), nil
}

func (bibm *Beis2Bums) BodyJson(header api.HeaderMap) ([]byte, error) {
	val := fastjson.Arena{}
	body := val.NewObject()
	val = fastjson.Arena{}
	head := val.NewObject()
	elements := bibm.root.ChildElements()
	for _, c := range elements {
		if ok := bibm.bodyTrancoder(c, head, body, header); ok {
			bibm.root.RemoveChild(c)
		}
	}

	bibm.OtherBody(bibm.root, body)

	val = fastjson.Arena{}
	resp := val.NewObject()
	resp.Set("head", head)
	resp.Set("body", body)
	data := resp.MarshalTo(nil)
	return data, nil
}

func (bibm *Beis2Bums) bodyTrancoder(e *etree.Element, head, body *fastjson.Value, header api.HeaderMap) bool {
	switch e.Tag {
	case "AppHead":
		bibm.BodyHead(e, head, "appHead", header)
	case "SysHead":
		bibm.BodyHead(e, head, "sysHead", header)
	case "LOCAL_HEAD":
		bibm.BodyHead(e, head, "localHead", header)
	case "details":
		bibm.DetailsV1(e, body)
	default:
		return false
	}
	return true
}

func (bibm *Beis2Bums) BodyHead(e *etree.Element, head *fastjson.Value, key string, header api.HeaderMap) {
	elements := e.ChildElements()
	for _, c := range elements {
		bibm.BodyHead(c, head, c.Tag, header)
	}
	if len(elements) == 0 {
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
				header.Set("SpanId", e.Text())
			}
			if strings.EqualFold(key, "traceid") {
				header.Set("TraceId", e.Text())
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

func (bibm *Beis2Bums) OtherBody(root *etree.Element, body *fastjson.Value) {
	elements := root.ChildElements()
	flag := make(map[string]bool)
	for _, e := range elements {
		if flag[e.Tag] {
			continue
		}
		flag[e.Tag] = true
		if es := root.SelectElements(e.Tag); len(es) != 1 {
			bodyKey := bibm.BodyKey(e.Tag)
			erena := fastjson.Arena{}
			array := erena.NewArray()
			bibm.VisitArray(root, array, es)
			body.Set(bodyKey, array)
		}
	}
	elements = root.ChildElements()
	for _, e := range elements {
		bodyKey := bibm.BodyKey(e.Tag)
		fval := bibm.VisitObject(e)
		body.Set(bodyKey, fval)
	}
}

func (bibm *Beis2Bums) DetailsV1(root *etree.Element, body *fastjson.Value) {
	val := fastjson.Arena{}
	details := val.NewObject()
	elements := root.ChildElements()
	if len(elements) != 0 {
		bibm.OtherBody(root, details)
		body.Set("details", details)
	} else if len(root.Text()) != 0 {
		body.Set("details", val.NewString(root.Text()))
	} else {
		log.DefaultContextLogger.Warnf(bibm.ctx, "the tag:%s of value is empty in details", root.Tag)
	}
}

func (bibm *Beis2Bums) DetailsV2(root *etree.Element, body *fastjson.Value) {
	val := fastjson.Arena{}
	details := val.NewObject()
	elements := root.ChildElements()
	flag := make(map[string]bool)
	for _, c := range elements {
		if flag[c.Tag] {
			continue
		}
		flag[c.Tag] = true
		bodyKey := bibm.BodyKey(c.Tag)
		erena := fastjson.Arena{}
		array := erena.NewArray()
		details.Set(bodyKey, array)
		elements := root.SelectElements(c.Tag)
		bibm.VisitArray(root, array, elements)
		details.Set(bodyKey, array)
	}

	if len(elements) != 0 {
		body.Set("details", details)
	} else if len(root.Text()) != 0 {
		body.Set("details", val.NewString(root.Text()))
	} else {
		log.DefaultContextLogger.Warnf(bibm.ctx, "the tag:%s of value is empty in details", root.Tag)
	}
}

// 保证数组顺序
func (bibm *Beis2Bums) VisitArray(parent *etree.Element, array *fastjson.Value, elements []*etree.Element) {
	for index, e := range elements {
		array.SetArrayItem(index, bibm.VisitObject(e))
		parent.RemoveChild(e)
	}
}

func (bibm *Beis2Bums) VisitObject(root *etree.Element) *fastjson.Value {
	erena := fastjson.Arena{}
	elements := root.ChildElements()
	if len(elements) == 0 {
		val := fastjson.Arena{}
		if len(root.Text()) == 0 {
			log.DefaultContextLogger.Warnf(bibm.ctx, "the tag:%s of value is empty in VisitObject", root.Tag)
		}
		return val.NewString(root.Text())
	}
	// slice
	object := erena.NewObject()
	flag := make(map[string]bool)
	for _, c := range elements {
		if flag[c.Tag] {
			continue
		}
		flag[c.Tag] = true
		if es := root.SelectElements(c.Tag); len(es) != 1 {
			erena := fastjson.Arena{}
			array := erena.NewArray()
			bibm.VisitArray(root, array, es)
			key := bibm.BodyKey(c.Tag)
			object.Set(key, array)
		}
	}
	// object
	es := root.ChildElements()
	for _, c := range es {
		if elements := root.SelectElements(c.Tag); len(elements) == 1 {
			key := bibm.BodyKey(c.Tag)
			object.Set(key, bibm.VisitObject(c))
		}
	}
	return object
}
