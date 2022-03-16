package main

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/protocol/http"
)

type BumsReq2BeisReq struct {
	header api.HeaderMap

	head    *fastjson.Value  // kv struct
	body    *fastjson.Object // list or string
	bodyErr error
	config  Bums2BeisConfig
	vo      Bums2BeisVo
}

func NewBumsReq2BeisReq(header api.HeaderMap, value string, config Bums2BeisConfig) (*BumsReq2BeisReq, error) {
	br2br := &BumsReq2BeisReq{
		config: config,
		header: header,
	}
	// json parse
	body, err := fastjson.Parse(value)
	if err != nil {
		return br2br, err
	}
	br2br.body, err = body.Get("body").Object()
	if err != nil {
		return br2br, err
	}

	br2br.head = body.Get("head")
	return br2br, nil
}

func (br2br *BumsReq2BeisReq) Head() (api.HeaderMap, error) {
	// 大小写转换
	beisHeader := http.RequestHeader{&fasthttp.RequestHeader{}}
	br2br.header.Range(func(k, v string) bool {
		k = strings.ToLower(k)
		beisHeader.Set(k, v)
		return true
	})

	origsender, ok := beisHeader.Get("origsender")
	if !ok || len(origsender) > 10 {
		return nil, fmt.Errorf("the %s of origsender is illage", origsender)
	}
	ctrlbits, ok := beisHeader.Get("ctrlbits")
	if !ok || len(ctrlbits) > 8 {
		return nil, fmt.Errorf("the %s of ctrlbits is illage", ctrlbits)
	}
	areacode, ok := beisHeader.Get("areacode")
	if !ok || len(areacode) != 4 {
		return nil, fmt.Errorf("the %s of areacode is illage", areacode)
	}
	versionid, ok := beisHeader.Get("versionid")
	if !ok || len(versionid) != 4 {
		return nil, fmt.Errorf("the %s of versionid is illage", versionid)
	}

	// 校验traceid ，spanid
	traceid, ok := beisHeader.Get("traceid")
	if ok {
		return nil, fmt.Errorf("the %s of traceid is illage", traceid)
	}

	spanid, ok := beisHeader.Get("spanid")
	if ok {
		return nil, fmt.Errorf("the %s of spanid is illage", spanid)
	}

	if len(br2br.vo.MesgId) != 0 {
		// beisHeader
		// TODO
	}

	if len(br2br.vo.MesgRefId) != 0 {
		// beisHeader
		// TODO
	}

	if len(br2br.vo.Reserve) != 0 {
		// beisHeader
		// TODO
	}
	return beisHeader, nil
}

func (br2br *BumsReq2BeisReq) CheckParam() bool {
	if len(br2br.vo.Namespace) != 20 {
		// TODO log
		return false
	} else if len(br2br.vo.MesgRefId) != 0 && len(br2br.vo.MesgRefId) > 20 {
		// TODO log
		return false
	} else if len(br2br.vo.Reserve) != 0 && len(br2br.vo.Reserve) > 45 {
		// TODO log
		return false
	}
	return true
}

func (br2br *BumsReq2BeisReq) GetXmlBytes(header api.HeaderMap) ([]byte, error) {
	beis := etree.NewDocument()
	beis.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := beis.CreateElement("Document")
	element.CreateAttr("xmlns", br2br.config.Namespace)

	sysHead := element.CreateElement("SysHead")
	if err := br2br.SysHead(sysHead); err != nil {
		return nil, err
	}

	appHead := element.CreateElement("AppHead")
	if err := br2br.AppHead(appHead, header); err != nil {
		return nil, err
	}
	// bugfix
	if err := br2br.Body(element); err != nil {
		return nil, err
	}
	return beis.WriteToBytes()
}

func (br2br *BumsReq2BeisReq) Transcoder() (api.HeaderMap, api.IoBuffer, error) {
	header, err := br2br.Head()
	if err != nil {
		return nil, nil, err
	}
	body, err := br2br.GetXmlBytes(header)
	if err != nil {
		return nil, nil, err
	}
	return header, buffer.NewIoBufferBytes(body), nil
}

func (br2br *BumsReq2BeisReq) SysHead(sysHead *etree.Element) error {
	for _, key := range br2br.config.SysHead {
		appKey := br2br.BumsHeadKey(key)
		val := br2br.head.GetStringBytes(appKey)
		if len(val) == 0 {
			return fmt.Errorf("the %s is not exist in sysHead", appKey)
		}
		sysHead.CreateElement(key).SetText(b2s(val))
	}
	return nil
}

func (br2br *BumsReq2BeisReq) AppHead(appHead *etree.Element, headers api.HeaderMap) error {
	traceid, _ := br2br.header.Get("traceid")
	appHead.CreateElement("Traceid").SetText(traceid)
	spanid, _ := br2br.header.Get("spanid")
	appHead.CreateElement("Spanid").SetText(spanid)

	for _, key := range br2br.config.AppHead {
		// 不区分大小写比较
		if strings.EqualFold(key, "traceid") || strings.EqualFold(key, "spanid") {
			continue
		}
		// 不区分大小写比较
		appKey := br2br.BumsHeadKey(key)
		if strings.EqualFold(key, "uniqueid") {
			val := br2br.head.GetStringBytes(appKey)
			if len(val) == 0 {
				return fmt.Errorf("the %s is not exist in appHead", appKey)
			}
			appHead.CreateElement("Uniqueid").SetText(b2s(val))
			continue
		}

		val := br2br.head.GetStringBytes(appKey)
		if len(val) == 0 {
			return fmt.Errorf("the %s is not exist in appHead", appKey)
		}
		appHead.CreateElement(key).SetText(b2s(val))
	}
	return nil
}

func (br2br *BumsReq2BeisReq) BumsHeadKey(key string) string {
	switch key {
	case "Uniqueid":
		return "uniqueId"
	case "Traceid":
		return "traceId"
	case "Spanid":
		return "spanId"
	}
	return ToFristLower(key)
}

func (br2br *BumsReq2BeisReq) Body(body *etree.Element) error {
	detailKey := "details"
	details := body.CreateElement(detailKey)
	conv := br2br.config.BodySwitch
	br2br.body.Visit(br2br.VisitRoot(body, details, conv))

	if br2br.bodyErr != nil {
		return br2br.bodyErr
	}

	detailsTokens := body.SelectElements(detailKey)
	for _, token := range detailsTokens {
		if len(token.Child) == 0 {
			body.RemoveChild(token)
		}
	}
	return nil
}

func (br2br *BumsReq2BeisReq) VisitRoot(body, details *etree.Element, conv bool) func([]byte, *fastjson.Value) {
	return func(key []byte, value *fastjson.Value) {
		if br2br.bodyErr != nil {
			return
		}
		switch value.Type() {
		case fastjson.TypeString:
			key = br2br.BodyKey(key, conv)
			v, _ := value.StringBytes()
			body.CreateElement(b2s(key)).SetText(b2s(v))
		case fastjson.TypeObject:
			key = br2br.BodyKey(key, conv)
			element := body.CreateElement(b2s(key))
			obj, _ := value.Object()
			obj.Visit(br2br.VisitList(element, conv))
		case fastjson.TypeArray:
			objs, _ := value.Array()
			key = br2br.BodyKey(key, br2br.config.DetailSwitch)
			for _, val := range objs {
				element := details.CreateElement(b2s(key))
				br2br.Array(element, val, key, br2br.config.DetailSwitch)
			}
		default:
			br2br.bodyErr = fmt.Errorf("not support:%s", value.Type().String())
		}
	}
}

func (br2br *BumsReq2BeisReq) VisitList(body *etree.Element, conv bool) func([]byte, *fastjson.Value) {
	return func(key []byte, value *fastjson.Value) {
		if br2br.bodyErr != nil {
			return
		}
		switch value.Type() {
		case fastjson.TypeString:
			key = br2br.BodyKey(key, conv)
			v, _ := value.StringBytes()
			body.CreateElement(b2s(key)).SetText(b2s(v))
		case fastjson.TypeObject:
			key = br2br.BodyKey(key, conv)
			element := body.CreateElement(b2s(key))
			obj, _ := value.Object()
			obj.Visit(br2br.VisitList(element, conv))
		case fastjson.TypeArray:
			objs, _ := value.Array()
			key = br2br.BodyKey(key, conv)
			for _, val := range objs {
				element := body.CreateElement(b2s(key))
				br2br.Array(element, val, key, conv)
			}
		default:
			br2br.bodyErr = fmt.Errorf("not support:%s", value.Type().String())
		}
	}
}

func (br2br *BumsReq2BeisReq) Array(body *etree.Element, value *fastjson.Value, key []byte, conv bool) {
	switch value.Type() {
	case fastjson.TypeString:
		v, _ := value.StringBytes()
		body.CreateElement(b2s(key)).SetText(b2s(v))
	case fastjson.TypeObject:
		obj, _ := value.Object()
		obj.Visit(br2br.VisitList(body, conv))
	case fastjson.TypeArray:
		objs, _ := value.Array()
		for _, val := range objs {
			br2br.Array(body, val, key, conv)
		}
	default:
		br2br.bodyErr = fmt.Errorf("not support:%s", value.Type().String())
	}
}

func (br2br *BumsReq2BeisReq) BodyKey(key []byte, conv bool) []byte {
	if conv {
		key = BytesToFristUpper(key)
	}
	return key
}
