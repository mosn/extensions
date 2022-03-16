package bumscd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/valyala/fastjson"
	"mosn.io/api"
)

const (
	HeadSource    = "head"
	HeadersSource = "headers"

	// Filed Support Type
	FieldString = "string"
	FieldImage  = "image"
	FieldByte   = "byte"
	FieldShort  = "short"
	FieldInt24  = "int24"
	FieldInt    = "int"
	FieldLong   = "long"
	FieldFloat  = "float"
	FieldDouble = "double"
	FieldList   = "list"
)

type BumsAndCdIterm struct {
	ListIterms map[string]*BumsAndCdIterm `json:"list_iterms,omitempty"`
	HeadIterms []*BumsAndCdIterm          `json:"head_iterms,omitempty"`
	CdKey      string                     `json:"cd_key,omitempty"`
	Type       string                     `json:"type,omitempty"`
	Length     string                     `json:"length,omitempty"`
	Scale      string                     `json:"scale,omitempty"`
	Describe   string                     `json:"describes,omitempty"`
	// head@tranTimestamp@
	Source  string `json:"source,omitempty"`
	BumsKey string `json:"bums_key,omitempty"`
	Default string `json:"default,omitempty"`
}

func NewBumsAndCdIterm(element *etree.Element, isResp bool) (*BumsAndCdIterm, error) {
	iterm := &BumsAndCdIterm{
		CdKey: element.Tag,
	}
	if err := iterm.HeadParse(element, isResp); err != nil {
		return nil, err
	}

	if strings.Contains(iterm.CdKey, "___") {
		iterm.CdKey = strings.Replace(iterm.CdKey, "___", "@", -1)
	}
	return iterm, nil
}

func NewBumsAndCdBodyIterm(element *etree.Element, isResp bool) (*BumsAndCdIterm, error) {
	var iterm = &BumsAndCdIterm{}
	if isResp {
		iterm.CdKey = element.Tag
	} else {
		iterm.BumsKey = element.Tag
	}
	if err := iterm.BodyParse(element, isResp); err != nil {
		return nil, err
	}
	if strings.Contains(iterm.CdKey, "___") {
		iterm.CdKey = strings.Replace(iterm.CdKey, "___", "@", -1)
	}
	return iterm, nil
}

func NewBumsAndCdList(element *etree.Element, list []*BumsAndCdIterm, isResp bool) (*BumsAndCdIterm, error) {
	baci := &BumsAndCdIterm{
		Type:       FieldList,
		HeadIterms: list,
		CdKey:      element.Tag,
	}
	if strings.Contains(baci.CdKey, "___") {
		baci.CdKey = strings.Replace(baci.CdKey, "___", "@", -1)
	}

	if isResp {
		attr, err := GetAttr(element, "convert")
		if err != nil || len(baci.CdKey) == 0 {
			baci.BumsKey = toSmallCamel(element.Tag)
		} else {
			baci.BumsKey = attr
		}
		return baci, nil
	}

	attr, err := GetAttr(element, "positionOrDefault")
	if err != nil {
		return nil, err
	}
	vals := strings.Split(attr, "@")
	if len(vals) < 3 {
		return nil, fmt.Errorf("the length of %s less three", attr)
	}
	baci.Source = vals[0]
	baci.BumsKey = vals[1]
	baci.Default = vals[2]
	return baci, nil
}

func NewBumsAndCdBodyList(element *etree.Element, list map[string]*BumsAndCdIterm, isResp bool) (*BumsAndCdIterm, error) {
	baci := &BumsAndCdIterm{
		ListIterms: list,
		Type:       FieldList,
	}
	var err error
	if isResp {
		baci.CdKey = element.Tag
		baci.BumsKey, err = GetAttr(element, "convert")
		if err != nil {
			return nil, err
		}
	} else {
		baci.BumsKey = element.Tag
		baci.CdKey, err = GetAttr(element, "convert")
		if err != nil {
			return nil, err
		}
	}
	if strings.Contains(baci.CdKey, "___") {
		baci.CdKey = strings.Replace(baci.CdKey, "___", "@", -1)
	}
	return baci, nil
}

func (baci *BumsAndCdIterm) HeadParse(element *etree.Element, isResp bool) (err error) {
	if isResp {
		baci.Type, err = GetAttr(element, "type")
		if err != nil {
			baci.Type = FieldString
		}
		baci.BumsKey, err = GetAttr(element, "convert")
		if err != nil || len(baci.BumsKey) == 0 {
			baci.BumsKey = toSmallCamel(element.Tag)
		}
		return nil
	}
	baci.Type, err = GetAttr(element, "type")
	if err != nil {
		return err
	}
	baci.Length, err = GetAttr(element, "length")
	if err != nil {
		return err
	}
	baci.Scale, err = GetAttr(element, "scale")
	if err != nil {
		return err
	}
	attr, err := GetAttr(element, "positionOrDefault")
	if err != nil {
		return err
	}
	vals := strings.Split(attr, "@")
	if len(vals) < 3 {
		return fmt.Errorf("the length of %s less three", attr)
	}
	baci.Source = vals[0]
	baci.BumsKey = vals[1]
	baci.Default = vals[2]
	return nil
}

func (baci *BumsAndCdIterm) BodyParse(element *etree.Element, isResp bool) (err error) {
	baci.Type, err = GetAttr(element, "type")
	if err != nil {
		return err
	}
	baci.Length, err = GetAttr(element, "length")
	if err != nil {
		return err
	}
	if isResp {
		baci.BumsKey, err = GetAttr(element, "convert")
		if err != nil {
			return err
		}

	} else {
		baci.CdKey, err = GetAttr(element, "convert")
		if err != nil {
			return err
		}
	}
	baci.Describe, err = GetAttr(element, "cnName")
	if err != nil {
		// TODO log
	}
	baci.Scale, err = GetAttr(element, "scale")
	if err != nil {
		baci.Scale = "0"
		// TODO log
	}
	return nil
}

func (baci *BumsAndCdIterm) GetValue(header api.HeaderMap, body *fastjson.Value) (string, error) {
	if len(baci.Source) == 0 {
		return baci.Default, nil
	}
	switch baci.Source {
	case HeadersSource:
		val, ok := header.Get(baci.BumsKey)
		if !ok {
			return "", fmt.Errorf("the %s of value is no exist", baci.BumsKey)
		}
		return val, nil
	case HeadSource:
		val := body.Get(baci.BumsKey)
		if val == nil {
			return "", fmt.Errorf("the %s of value is no exist", baci.BumsKey)
		}
		return baci.GetBodyValue(val)
	default:
		return "", fmt.Errorf("the %s of type is support", baci.Source)
	}
}

func (baci *BumsAndCdIterm) GetBodyValue(body *fastjson.Value) (string, error) {
	value, err := GetValue(baci.Type, body)
	if err != nil {
		return "", err
	}
	return B2S(value), nil
}

func (baci *BumsAndCdIterm) String() string {
	v, _ := json.Marshal(baci)
	return B2S(v)
}

func (baci *BumsAndCdIterm) CheckParam() error {
	if baci.Source != "" && baci.Source != HeadersSource && baci.Source != HeadSource {
		return fmt.Errorf("the source of %s is support", baci.Source)
	}
	return nil
}

func GetValue(itype string, fval *fastjson.Value) (value []byte, err error) {
	switch fval.Type() {
	case fastjson.TypeString:
		value, err = fval.StringBytes()
	case fastjson.TypeNumber:
		return GetNumber(itype, fval)
	default:
		return nil, fmt.Errorf("the type of %s support", fval.Type())
	}
	return value, err
}

func GetNumber(itype string, fval *fastjson.Value) (value []byte, err error) {
	switch itype {
	case FloatField:
		fallthrough
	case DoubleField:
		fdata, err := fval.Float64()
		if err != nil {
			return nil, err
		}
		sdata := strconv.FormatFloat(fdata, 'g', -1, 64)
		return S2B(sdata), nil
	default:
		val, err := fval.Int()
		if err != nil {
			return nil, err
		}
		return S2B(strconv.Itoa(val)), nil
	}
}

func GetAttr(element *etree.Element, key string) (string, error) {
	attr := element.SelectAttr(key)
	if attr == nil {
		return "", fmt.Errorf("the %s of select attr is empty", key)
	}
	return attr.Value, nil
}

func GetObject(iterm BumsAndCdIterm, value string) (interface{}, error) {
	switch iterm.Type {
	case FieldImage:
		fallthrough
	case ByteField:
		fallthrough
	case StringField:
		return value, nil
	case ShortField:
		fallthrough
	case Int24Field:
		fallthrough
	case IntField:
		fallthrough
	case LongField:
		ival, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return ival, nil
	case DoubleField:
		fallthrough
	case FloatField:
		fval, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return fval, nil
	default:
		return value, nil
	}
}
