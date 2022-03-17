package bumsbeis

import (
	"context"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

type Bums2Beis struct {
	ctx    context.Context
	header api.HeaderMap

	head    *fastjson.Value  // kv struct
	body    *fastjson.Object // list or string
	bodyErr error
	config  *Bums2BeisConfig
	vo      *Bums2BeisVo
}

func NewBums2Beis(ctx context.Context, header api.HeaderMap, buf api.IoBuffer, config *Bums2BeisConfig, vo *Bums2BeisVo) (*Bums2Beis, error) {
	br2br := &Bums2Beis{
		config: config,
		vo:     vo,
		header: header,
	}
	// json parse
	body, err := fastjson.Parse(buf.String())
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

func (br2br *Bums2Beis) CheckParam() error {
	if len(br2br.vo.Namespace) != 20 {
		return fmt.Errorf("the %s of namespace is illage", br2br.vo.Namespace)
	} else if len(br2br.vo.MesgId) != 0 && len(br2br.vo.MesgRefId) > 20 {
		return fmt.Errorf("the %s of mesgId is illage", br2br.vo.MesgId)
	} else if len(br2br.vo.MesgRefId) != 0 && len(br2br.vo.MesgRefId) > 20 {
		return fmt.Errorf("the %s of mesgRefId is illage", br2br.vo.MesgRefId)
	} else if len(br2br.vo.Reserve) != 0 && len(br2br.vo.Reserve) > 45 {
		return fmt.Errorf("the %s of reserve is illage", br2br.vo.Reserve)
	}
	return nil
}

func (br2br *Bums2Beis) GetXmlBytes(header api.HeaderMap) ([]byte, error) {
	beis := etree.NewDocument()
	beis.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := beis.CreateElement("Document")
	element.CreateAttr("xmlns", br2br.vo.Namespace)

	sysHead := element.CreateElement("SysHead")
	if err := br2br.SysHead(sysHead, header); err != nil {
		return nil, err
	}

	appHead := element.CreateElement("AppHead")
	if err := br2br.AppHead(appHead, header); err != nil {
		return nil, err
	}

	if err := br2br.Body(element); err != nil {
		return nil, err
	}
	return beis.WriteToBytes()
}

func (br2br *Bums2Beis) Transcoder(isRequest bool) (header api.HeaderMap, buf api.IoBuffer, err error) {
	if isRequest {
		header, err = br2br.HeadRequest()
	} else {
		header, err = br2br.HeadRespone()
	}
	if err != nil {
		return nil, nil, err
	}
	body, err := br2br.GetXmlBytes(header)
	if err != nil {
		return nil, nil, err
	}
	return header, buffer.NewIoBufferBytes(body), nil
}

func (br2br *Bums2Beis) SysHead(sysHead *etree.Element, header api.HeaderMap) error {
	for _, key := range br2br.config.SysHead {
		appKey := br2br.BumsHeadKey(key)
		val := br2br.head.GetStringBytes(appKey)
		if len(val) == 0 {
			return fmt.Errorf("the %s is not exist in sysHead", appKey)
		}
		sysHead.CreateElement(key).SetText(b2s(val))

		// 逻辑兼容
		if strings.EqualFold(key, "ServiceScene") {
			header.Set("ServiceScene", b2s(val))
		}
		if strings.EqualFold(key, "ServiceCode") {
			header.Set("ServiceCode", b2s(val))
		}
	}
	return nil
}

func (br2br *Bums2Beis) AppHead(appHead *etree.Element, headers api.HeaderMap) error {
	traceid, _ := headers.Get("traceid")
	appHead.CreateElement("Traceid").SetText(traceid)
	spanid, _ := headers.Get("spanid")
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

func (br2br *Bums2Beis) BumsHeadKey(key string) string {
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

func (br2br *Bums2Beis) Body(body *etree.Element) error {
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

func (br2br *Bums2Beis) VisitRoot(body, details *etree.Element, conv bool) func([]byte, *fastjson.Value) {
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

func (br2br *Bums2Beis) VisitList(body *etree.Element, conv bool) func([]byte, *fastjson.Value) {
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

func (br2br *Bums2Beis) Array(body *etree.Element, value *fastjson.Value, key []byte, conv bool) {
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

func (br2br *Bums2Beis) BodyKey(key []byte, conv bool) []byte {
	if conv {
		key = BytesToFristUpper(key)
	}
	return key
}
