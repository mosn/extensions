package bumsbeis

import (
	"strings"
	"unsafe"
)

type Bums2BeisVo struct {
	Namespace string `json:"namespace"`
	MesgId    string `json:"mesgId"`
	MesgRefId string `json:"mesg_refId"`
	Reserve   string `json:"reserve"`
}

type Bums2BeisConfig struct {
	SysHead      []string `json:"sys_head"`
	AppHead      []string `json:"app_head"`
	DetailSwitch bool     `json:"detail_switch"`
	BodySwitch   bool     `json:"body_switch"`
	//
	Namespace string `json:"namespace"`
}

// not support utf8.RuneSelf
func ToFristLower(r string) string {
	if len(r) <= 0 {
		return r
	}
	if r[0] >= 'a' && r[0] <= 'z' {
		return r
	}
	var b strings.Builder
	b.Grow(len(r))
	c := r[0] + 'a' - 'A'
	b.WriteByte(c)
	b.WriteString(r[1:])
	return b.String()
}

func ToFristUpper(r string) string {
	if len(r) <= 0 {
		return r
	}
	if r[0] >= 'A' && r[0] <= 'Z' {
		return r
	}
	var b strings.Builder
	b.Grow(len(r))
	c := r[0] - 'a' + 'A'
	b.WriteByte(c)
	b.WriteString(r[1:])
	return b.String()
}

// not support utf8.RuneSelf
func BytesToFristUpper(r []byte) []byte {
	if len(r) <= 0 || (r[0] >= 'A' && r[0] <= 'Z') {
		return r
	}
	r[0] -= 'a' - 'A'
	return r
}

func BytesToFristLower(r []byte) []byte {
	if len(r) <= 0 || (r[0] >= 'a' && r[0] <= 'z') {
		return r
	}
	r[0] += 'a' - 'A'
	return r
}

// TODO b panic
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
