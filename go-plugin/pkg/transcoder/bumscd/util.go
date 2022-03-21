package bumscd

import (
	"reflect"
	"strings"
	"unsafe"
)

const (
	// xml key
	DataXML   = "data"
	NameXML   = "name"
	FieldXML  = "field"
	ArrayXML  = "array"
	LengthXML = "length"
	ScaleXML  = "scale"
	TypeXML   = "type"
	StructXML = "struct"

	StringField = "string"
	ImageField  = "image"
	ByteField   = "byte"
	ShortField  = "short"
	Int24Field  = "int24"
	IntField    = "int"
	LongField   = "long"
	FloatField  = "float"
	DoubleField = "double"
	ListField   = "list"
)

func B2S(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

func S2B(s string) (b []byte) {
	strh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Data = strh.Data
	sh.Len = strh.Len
	sh.Cap = strh.Len
	return b
}

func toSmallCamel(name string) string {
	sb := strings.Builder{}
	tempName := strings.Split(strings.ToLower(name), "_")
	for i := 0; i < len(tempName); i++ {
		if i != 0 {
			chars := []byte(tempName[i])
			c := (chars[0] - 32)
			sb.WriteByte(c)
			sb.Write(chars[1:])
		} else {
			sb.WriteString(tempName[0])
		}
	}
	return sb.String()
}
