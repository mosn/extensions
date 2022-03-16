package bumscd

import (
	"fmt"
	"strconv"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
	"mosn.io/api"
)

type Cd2Bums struct {
	root   *etree.Element
	config *Relation
	header api.HeaderMap
}

func NewCd2Bums(header api.HeaderMap, buf api.IoBuffer, config *Relation) (*Cd2Bums, error) {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(buf)
	if err != nil {
		return nil, err
	}
	root := doc.SelectElement("service")
	if root == nil {
		return nil, fmt.Errorf("doc is empty:doc :%s", doc.FullTag())
	}
	return &Cd2Bums{
		root:   root,
		header: header,
		config: config,
	}, nil
}

func (cd2bm *Cd2Bums) Body() (*fastjson.Value, error) {
	arena := fastjson.Arena{}
	obj := arena.NewObject()

	head, err := cd2bm.HeadInBody()
	if err != nil {
		return nil, err
	}
	obj.Set("head", head)

	body, err := cd2bm.BodyInBody()
	if err != nil {
		return nil, err
	}
	obj.Set("body", body)
	return obj, nil
}

// 数据排平 写成kv结构，类型保持string:string
func (cd2bm *Cd2Bums) HeadInBody() (*fastjson.Value, error) {
	arena := fastjson.Arena{}
	obj := arena.NewObject()

	key := "sys-header"
	if err := cd2bm.bodyHead(key, obj, cd2bm.config.SysHead); err != nil {
		return nil, err
	}

	key = "app-header"
	if err := cd2bm.bodyHead(key, obj, cd2bm.config.AppHead); err != nil {
		return nil, err
	}

	key = "local-header"
	if err := cd2bm.bodyHead(key, obj, cd2bm.config.LocalHead); err != nil {
		return nil, err
	}
	return obj, nil
}

func (cd2bm *Cd2Bums) bodyHead(key string, obj *fastjson.Value, configs []*BumsAndCdIterm) error {
	head := cd2bm.root.SelectElement(key)
	if head == nil {
		return fmt.Errorf("the %s of head is empty in bodyhead", key)
	}
	data := head.SelectElement(DataXML)
	if data == nil {
		return fmt.Errorf("the %v of elements is empty in bodyhead", head)
	}
	child := data.SelectElement(StructXML)
	if child == nil {
		return fmt.Errorf("the %v of element is empty in bodyhead", data)
	}
	return cd2bm.writeHeadData(child, configs, obj)
}

func (cd2bm *Cd2Bums) writeHeadData(root *etree.Element, configs []*BumsAndCdIterm, obj *fastjson.Value) error {
	for _, esub := range root.ChildElements() {
		config, err := cd2bm.selectConfig(esub, configs)
		if err != nil {
			// TODO log
			continue
		}
		for _, e := range esub.ChildElements() {
			switch e.Tag {
			case "field":
				fval, err := cd2bm.parseField(e, config)
				if err != nil {
					return err
				}
				obj.Set(config.BumsKey, fval)
			case "array":
				if config.Type != FieldList {
					return fmt.Errorf("the %s is illage in array", config.Type)
				}
				if config.CdKey == "RET" {
					if err := cd2bm.writeRet(e, config.HeadIterms, obj); err != nil {
						return err
					}
				} else {
					farray, err := cd2bm.writeArray(e, config.HeadIterms)
					if err != nil {
						return err
					}
					obj.Set(config.BumsKey, farray)
				}
			default:
				return fmt.Errorf("the %s of tag not support", e.Tag)
			}
		}
	}
	return nil
}

func (cd2bm *Cd2Bums) selectConfig(root *etree.Element, config []*BumsAndCdIterm) (*BumsAndCdIterm, error) {
	wname := root.SelectAttr(NameXML)
	if wname == nil {
		return nil, fmt.Errorf("the %s of element is attr in data", NameXML)
	}
	for _, c := range config {
		if c.CdKey == wname.Value {
			return c, nil
		}
	}
	// Mock
	return &BumsAndCdIterm{BumsKey: toSmallCamel(wname.Value), Type: StringField, CdKey: wname.Value}, nil
}

func (cd2bm *Cd2Bums) writeArray(root *etree.Element, configs []*BumsAndCdIterm) (*fastjson.Value, error) {
	elements := root.ChildElements()
	arena := fastjson.Arena{}
	farray := arena.NewArray()
	for index, e := range elements {
		switch e.Tag {
		case "array":
			obj, err := cd2bm.writeArray(e, configs)
			if err != nil {
				return nil, err
			}
			farray.SetArrayItem(index, obj)
		case "struct":
			arena := fastjson.Arena{}
			obj := arena.NewObject()
			farray.SetArrayItem(index, obj)
			if err := cd2bm.writeHeadData(e, configs, obj); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("the %s of tag is not support in array", e.Tag)
		}
	}
	return farray, nil
}

func (cd2bm *Cd2Bums) parseField(root *etree.Element, config *BumsAndCdIterm) (*fastjson.Value, error) {
	rtype := root.SelectAttr(TypeXML)
	if rtype == nil {
		return nil, fmt.Errorf("the %s of attr is not exist", TypeXML)
	}
	switch rtype.Value {
	case FieldByte:
		fallthrough
	case FieldString:
		arena := fastjson.Arena{}
		return arena.NewString(root.Text()), nil
	case FieldShort:
		fallthrough
	case FieldInt24:
		fallthrough
	case FieldInt:
		fallthrough
	case FieldLong:
		value := root.Text()
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		arena := fastjson.Arena{}
		return arena.NewNumberInt(num), nil
	case FieldFloat:
		fallthrough
	case FieldDouble:
		value := root.Text()
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		arena := fastjson.Arena{}
		return arena.NewNumberFloat64(num), nil
	default:
		return nil, fmt.Errorf("the %s of type not support", config.Type)
	}
}

// 特殊字段处理
func (cd2bm *Cd2Bums) writeRet(root *etree.Element, config []*BumsAndCdIterm, obj *fastjson.Value) error {
	array, err := cd2bm.parseStruct(root)
	if err != nil {
		return err
	}

	elements := array.ChildElements()
	for _, data := range elements {
		name := data.SelectAttr(NameXML)
		if name == nil {
			return fmt.Errorf("the %v of name is empty in writeRet", data)
		}
		switch name.Value {
		case "RET_CODE":
			field := data.SelectElement(FieldXML)
			if field == nil {
				// TODO log
				break
			}
			if field.Text() != "" {
				tmp := fastjson.Arena{}
				val := tmp.NewString(field.Text())
				obj.Set("retCode", val)
			}
		case "RET_MSG":
			field := data.SelectElement(FieldXML)
			if field == nil {
				// TODO log
				break
			}
			if field.Text() != "" {
				tmp := fastjson.Arena{}
				val := tmp.NewString(field.Text())
				obj.Set("retMsg", val)
			}
		}
	}
	return nil
}

func (cd2bm *Cd2Bums) parseStruct(root *etree.Element) (*etree.Element, error) {
	child := root.SelectElement(StructXML)
	if child == nil {
		return nil, fmt.Errorf("the %v of element is empty in struct", root)
	}
	return child, nil
}

func (cd2bm *Cd2Bums) BodyInBody() (*fastjson.Value, error) {
	arena := fastjson.Arena{}
	obj := arena.NewObject()
	body := cd2bm.root.SelectElement("body")
	if body == nil {
		return nil, fmt.Errorf("the %v of element is empty in body", cd2bm.root)
	}
	if err := cd2bm.writeBodyData(body, cd2bm.config.Body, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (cd2bm *Cd2Bums) writeBodyData(root *etree.Element, configs map[string]*BumsAndCdIterm, obj *fastjson.Value) error {
	for _, esub := range root.ChildElements() {
		wname := esub.SelectAttr(NameXML)
		if wname == nil {
			return fmt.Errorf("the %s of element is attr in bodydata", NameXML)
		}
		config, ok := configs[wname.Value]
		if !ok {
			// TODO log
			continue
		}
		for _, e := range esub.ChildElements() {
			switch e.Tag {
			case "field":
				fval, err := cd2bm.parseField(e, config)
				if err != nil {
					return err
				}
				obj.Set(config.BumsKey, fval)
			case "array":
				if config.Type != FieldList {
					return fmt.Errorf("the %s is illage in array", config.Type)
				}
				farray, err := cd2bm.writeBodyArray(e, config.ListIterms)
				if err != nil {
					return err
				}
				obj.Set(config.BumsKey, farray)
			case "struct":
				arena := fastjson.Arena{}
				objn := arena.NewObject()
				if err := cd2bm.writeBodyData(e, configs, objn); err != nil {
					return err
				}
				obj.Set(config.BumsKey, objn)
			default:
				return fmt.Errorf("the %s of tag not support", e.Tag)
			}
		}
	}
	return nil
}

func (cd2bm *Cd2Bums) writeBodyArray(root *etree.Element, configs map[string]*BumsAndCdIterm) (*fastjson.Value, error) {
	elements := root.ChildElements()
	arena := fastjson.Arena{}
	farray := arena.NewArray()
	for index, e := range elements {
		switch e.Tag {
		case "array":
			wname := e.SelectAttr(NameXML)
			if wname == nil {
				return nil, fmt.Errorf("the %s of element is attr in data", NameXML)
			}
			config, ok := configs[wname.Value]
			if !ok {
				return nil, fmt.Errorf("the %s of config is empty in bodyarray", e.Tag)
			}
			obj, err := cd2bm.writeBodyArray(e, config.ListIterms)
			if err != nil {
				return nil, err
			}
			farray.SetArrayItem(index, obj)
		case "struct":
			arena := fastjson.Arena{}
			obj := arena.NewObject()
			farray.SetArrayItem(index, obj)
			if err := cd2bm.writeBodyData(e, configs, obj); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("the %s of tag is not support in array", e.Tag)
		}
	}
	return farray, nil
}
