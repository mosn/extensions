package common

import (
	"bytes"
	"github.com/mosn/extensions/go-plugin/pkg/protocol/dubbo/constants"
)

func BuildDubboDataId(name, version, group string) string {
	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	buffer.WriteString(name)
	if version != "" {
		buffer.WriteString(":")
		buffer.WriteString(version)
	}
	if group != "" {
		buffer.WriteString(":")
		buffer.WriteString(group)
	}
	buffer.WriteString(constants.XPROTOCOL_TYPE_DUBBO_DATAID_SUFFIX)
	return buffer.String()
}
