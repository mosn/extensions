package main

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

type BeisResp2BumsResp struct {
	root   *etree.Element
	header api.HeaderMap
}

func NewBeisResp2BumsResp(header api.HeaderMap, buf api.IoBuffer) (*BeisResp2BumsResp, error) {
	doc := etree.NewDocument()
	doc.ReadFrom(buf)
	root := doc.SelectElement("Document")
	if root == nil {
		return nil, fmt.Errorf("doc is empty:doc :%s", doc.FullTag())
	}
	return &BeisResp2BumsResp{
		root:   root,
		header: header,
	}, nil
}

func (br2br *BeisResp2BumsResp) BodyJson() ([]byte, error) {
	val := fastjson.Arena{}
	body := val.NewObject()
	val = fastjson.Arena{}
	head := val.NewObject()
	for _, t := range br2br.root.Child {
		if c, ok := t.(*etree.Element); ok {
			br2br.bodyTrancoder(c, head, body)
		}
	}

	val = fastjson.Arena{}
	resp := val.NewObject()
	resp.Set("head", head)
	resp.Set("body", body)
	data := resp.MarshalTo(nil)
	return data, nil
}

func (br2br *BeisResp2BumsResp) Transcoder() (api.HeaderMap, api.IoBuffer, error) {
	header := br2br.header
	body, err := br2br.BodyJson()
	if err != nil {
		return nil, nil, err
	}
	return header, buffer.NewIoBufferBytes(body), nil
}

func (br2br *BeisResp2BumsResp) bodyTrancoder(e *etree.Element, head, body *fastjson.Value) {
	switch e.Tag {
	case "AppHead":
		br2br.Head(e, head, "appHead")
	case "SysHead":
		br2br.Head(e, head, "sysHead")
	case "LOCAL_HEAD":
		br2br.Head(e, head, "localHead")
	case "details":
		br2br.Details(e, body)
	default:
		br2br.OtherBody(e, body, e.Tag, 0)
	}
}

func (br2br *BeisResp2BumsResp) Head(e *etree.Element, head *fastjson.Value, key string) {
	elemEmpty := true
	for _, t := range e.Child {
		if c, ok := t.(*etree.Element); ok {
			br2br.Head(c, head, c.Tag)
		}
	}
	if elemEmpty {
		if len(e.Text()) == 0 {
			// return fmt.Errorf("the %s of text is empty", e.Tag)
			return
		}
		val := fastjson.Arena{}
		head.Set(br2br.HeadKey(key), val.NewString(e.Text()))
	}
	// TODO bug
}

func (br2br *BeisResp2BumsResp) HeadKey(key string) string {
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

func (br2br *BeisResp2BumsResp) BodyKey(key string) string {
	return key
	// return ToFristLower(key)
}

func (br2br *BeisResp2BumsResp) OtherBody(e *etree.Element, body *fastjson.Value, key string, depth int) {
	elemEmpty := true
	val := fastjson.Arena{}
	obj := val.NewObject()

	for _, t := range e.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = false
			if depth <= 0 {
				br2br.OtherBody(c, obj, c.Tag, depth+1)
			} else {
				br2br.Head(c, obj, c.Tag)
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

func (br2br *BeisResp2BumsResp) Details(e *etree.Element, body *fastjson.Value) {
	val := fastjson.Arena{}
	details := val.NewObject()
	elemEmpty := true
	for _, t := range e.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = false
			bodyKey := br2br.BodyKey(c.Tag)
			array := details.Get(bodyKey)
			if array == nil {
				erena := fastjson.Arena{}
				array = erena.NewArray()
				details.Set(bodyKey, array)
			}
			elements := e.SelectElements(c.Tag)
			br2br.VisitArray(e, array, elements)
			details.Set(bodyKey, array)
		}
	}
	if !elemEmpty {
		body.Set("details", details)
	} else if len(e.Text()) != 0 {
		body.Set("details", val.NewString(e.Text()))
	}
	// TODO log
	// return fmt.Errorf("the %v of details is empty", e)
}

// 保证数组顺序
func (br2br *BeisResp2BumsResp) VisitArray(parent *etree.Element, array *fastjson.Value, elements []*etree.Element) {
	for index, e := range elements {
		array.SetArrayItem(index, br2br.VisitObject(e))
		parent.RemoveChild(e)
	}
}

func (br2br *BeisResp2BumsResp) VisitObject(ele *etree.Element) *fastjson.Value {
	erena := fastjson.Arena{}
	object := erena.NewObject()
	elemEmpty := true
	for _, t := range ele.Child {
		if c, ok := t.(*etree.Element); ok {
			elemEmpty = false
			if elements := ele.SelectElements(c.Tag); len(elements) == 1 {
				key := br2br.BodyKey(c.Tag)
				object.Set(key, br2br.VisitObject(c))
			} else {
				erena := fastjson.Arena{}
				array := erena.NewArray()
				br2br.VisitArray(ele, array, elements)
				key := br2br.BodyKey(c.Tag)
				object.Set(key, array)
			}
		}
	}
	if elemEmpty {
		val := fastjson.Arena{}
		if len(ele.Text()) == 0 {
			// TODO log
		}
		return val.NewString(ele.Text())
	}
	return object
}
