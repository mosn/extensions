package bumscd

import (
	"context"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

type Bums2Cd struct {
	ctx          context.Context
	header       api.HeaderMap
	config       *Config
	head         *fastjson.Value  // kv struct
	body         *fastjson.Object // list or string
	bodyVisitErr error
}

func NewBums2Cd(ctx context.Context, header api.HeaderMap, buf api.IoBuffer, config *Config) (*Bums2Cd, error) {
	bm2cd := &Bums2Cd{
		config: config,
		header: header,
		ctx:    ctx,
	}
	// json parse
	body, err := fastjson.Parse(buf.String())
	if err != nil {
		return nil, err
	}
	bm2cd.body, err = body.Get("body").Object()
	if err != nil {
		return nil, err
	}
	bm2cd.head = body.Get("head")
	if bm2cd.head == nil {
		return nil, fmt.Errorf("the %s of head is not exist", body)
	}
	return bm2cd, nil
}

func (bm2cd *Bums2Cd) Transcoder(isReq bool) (headers api.HeaderMap, buf api.IoBuffer, err error) {
	var relation *Relation
	if isReq {
		headers, err = bm2cd.HeadRequest()
		relation = bm2cd.config.ReqMapping
	} else {
		headers, err = bm2cd.HeadRespone()
		relation = bm2cd.config.RespMapping
	}
	if err != nil {
		return nil, nil, err
	}
	body, err := bm2cd.GetXmlBytes(relation)
	if err != nil {
		return nil, nil, err
	}
	return headers, buffer.NewIoBufferBytes(body), nil
}

func (bm2cd *Bums2Cd) GetXmlBytes(relation *Relation) ([]byte, error) {
	cd := etree.NewDocument()
	cd.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := cd.CreateElement("service")
	if err := bm2cd.BodyHead(element, relation); err != nil {
		return nil, err
	}
	if err := bm2cd.Body(element, relation); err != nil {
		return nil, err
	}
	return cd.WriteToBytes()
}

func (bm2cd *Bums2Cd) BodyHead(root *etree.Element, relation *Relation) error {
	sysHead := bm2cd.genHeadElement(root, "sys-header", "SYS_HEAD")
	bm2cd.writeHead(sysHead, relation.SysHead, bm2cd.header, bm2cd.head)
	appHead := bm2cd.genHeadElement(root, "app-header", "APP_HEAD")
	bm2cd.writeHead(appHead, relation.AppHead, bm2cd.header, bm2cd.head)
	localHead := bm2cd.genHeadElement(root, "local-header", "LOCAL_HEAD")
	bm2cd.writeHead(localHead, relation.LocalHead, bm2cd.header, bm2cd.head)
	return nil
}

func (bm2cd *Bums2Cd) genHeadElement(root *etree.Element, key, name string) *etree.Element {
	element := root.CreateElement(key)
	element = element.CreateElement(DataXML)
	element.CreateAttr(NameXML, name)
	return element
}

func (bm2cd *Bums2Cd) writeHead(root *etree.Element, headConfig []*BumsAndCdIterm, header api.HeaderMap, body *fastjson.Value) (err error) {
	var child *etree.Element
	child = root.CreateElement(StructXML)
	for _, v := range headConfig {
		switch v.Type {
		case FieldList:
			array := child.CreateElement(DataXML)
			array.CreateAttr(NameXML, v.CdKey)
			err = bm2cd.writeArray(array, v, header, body)
		default:
			err = bm2cd.writeFiled(child, v, header, body)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (bm2cd *Bums2Cd) writeArray(root *etree.Element, list *BumsAndCdIterm, header api.HeaderMap, body *fastjson.Value) (err error) {
	if len(list.HeadIterms) == 0 {
		return fmt.Errorf("the %s of array is empty", list.CdKey)
	}
	child := root.CreateElement(ArrayXML).CreateElement(StructXML)
	farray, err := body.Get(list.BumsKey).Array()
	if err != nil {
		return fmt.Errorf("the %s of body is empty", list.BumsKey)
	}
	for _, v := range list.HeadIterms {
		switch v.Type {
		case FieldList:
			return fmt.Errorf("the %s of type is not support", v.Type)
		default:
			// 源码逻辑兼容
			if subs := strings.Split(v.BumsKey, "#"); len(subs) > 1 {
				v.BumsKey = subs[1]
			}
			value, err := bm2cd.findValueInArray(v, farray)
			if err != nil {
				return err
			}
			if err := bm2cd.writeFiled(child, v, header, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bm2cd *Bums2Cd) findValueInArray(list *BumsAndCdIterm, array []*fastjson.Value) (*fastjson.Value, error) {
	for _, arr := range array {
		if val := arr.Get(list.BumsKey); val != nil {
			return arr, nil
		}
	}
	return nil, fmt.Errorf("the %s of value not exist", list.BumsKey)
}

func (bm2cd *Bums2Cd) writeFiled(root *etree.Element, iterm *BumsAndCdIterm, header api.HeaderMap, body *fastjson.Value) error {
	if iterm.Type == FieldList {
		return fmt.Errorf("the %s of type illage", iterm.Type)
	}
	field, err := bm2cd.filedElement(iterm, header, body)
	if err != nil {
		return err
	}
	element := etree.NewElement(DataXML)
	element.CreateAttr(NameXML, iterm.CdKey)
	element.AddChild(field)
	root.AddChild(element)
	// TODO code 处理
	return nil
}

func (bm2cd *Bums2Cd) filedElement(iterm *BumsAndCdIterm, header api.HeaderMap, body *fastjson.Value) (*etree.Element, error) {
	var value string
	var err error
	if header == nil {
		value, err = iterm.GetBodyValue(body)
	} else {
		value, err = iterm.GetValue(header, body)
	}
	if err != nil {
		return nil, err
	}

	e := etree.NewElement(FieldXML)
	e.SetText(value)
	e.CreateAttr(LengthXML, iterm.Length)
	e.CreateAttr(ScaleXML, iterm.Scale)
	e.CreateAttr(TypeXML, iterm.Type)
	return e, nil
}

func (bm2cd *Bums2Cd) Body(root *etree.Element, relation *Relation) error {
	body := root.CreateElement("body")
	bm2cd.body.Visit(bm2cd.visitBody(body, relation.Body))
	if bm2cd.bodyVisitErr != nil {
		return bm2cd.bodyVisitErr
	}
	return nil
}

func (bm2cd *Bums2Cd) visitBody(root *etree.Element, bodyConfig map[string]*BumsAndCdIterm) func([]byte, *fastjson.Value) {
	return func(key []byte, value *fastjson.Value) {
		if bm2cd.bodyVisitErr != nil {
			return
		}

		iterm, ok := bodyConfig[B2S(key)]
		if !ok {
			bm2cd.bodyVisitErr = fmt.Errorf("the %s of key is not exist", key)
			return
		}
		switch value.Type() {
		case fastjson.TypeNumber:
			fallthrough
		case fastjson.TypeString:
			if err := bm2cd.writeFiled(root, iterm, nil, value); err != nil {
				bm2cd.bodyVisitErr = err
				return
			}
		case fastjson.TypeArray:
			if err := bm2cd.visitArray(root, iterm, value); err != nil {
				bm2cd.bodyVisitErr = err
				return
			}
		default:
			bm2cd.bodyVisitErr = fmt.Errorf("the %s of type is not support", value.Type())
		}
	}
}

func (bm2cd *Bums2Cd) visitArray(root *etree.Element, list *BumsAndCdIterm, body *fastjson.Value) error {
	if list.Type != FieldList {
		return fmt.Errorf("the %s of type illage", list.Type)
	}
	array := root.CreateElement(DataXML)
	array.CreateAttr(NameXML, list.CdKey)
	child := array.CreateElement(ArrayXML)
	values, _ := body.Array()
	for _, value := range values {
		switch value.Type() {
		case fastjson.TypeObject:
			obj, _ := value.Object()
			obj.Visit(bm2cd.visitBody(child.CreateElement(StructXML), list.ListIterms))
			if bm2cd.bodyVisitErr != nil {
				return bm2cd.bodyVisitErr
			}
		default:
			return fmt.Errorf("not support:%s", value.Type())
		}
	}
	return nil
}
