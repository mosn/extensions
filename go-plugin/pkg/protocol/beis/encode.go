/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package beis

import (
	"bytes"
	"context"
	"encoding/xml"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
)

func (proto *Protocol) encodeRequest(ctx context.Context, request *Request) (api.IoBuffer, error) {

	packetLen := RequestHeaderLen /** fixed 128 byte length */ + request.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 128 byte length, 8 byte fixed begin flag
	buf.WriteString(beginFlag)
	// 10 byte origin sender
	proto.suffixOfBlank(buf, request.OrigSender, 10)
	// 8 byte message length
	proto.prefixOfZero(buf, request.Payload.Len(), 8)
	// 8 byte control bits.
	buf.WriteString(request.CtrlBits)
	// 4 byte AreaCode
	proto.suffixOfBlank(buf, request.AreaCode, 4)
	// 4 byte fixed version
	buf.WriteString("0001")
	// 20 byte
	proto.suffixOfBlank(buf, request.MessageID, 20)
	// 20 byte
	proto.suffixOfBlank(buf, request.MessageRefID, 20)
	// 45 byte
	proto.suffixOfBlank(buf, request.Reserve, 45)
	// 1 byte end flag
	buf.WriteString("}")

	if request.Payload.Len() > 0 {
		// 2. write payload bytes
		buf.Write(request.Payload.Bytes())
	}

	return encrypt(ctx, buf, packetLen)
}

func (proto *Protocol) encodeResponse(ctx context.Context, response *Response) (api.IoBuffer, error) {

	packetLen := RequestHeaderLen /** fixed 128 byte length */ + response.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 128 byte length, 8 byte fixed begin flag
	buf.WriteString(beginFlag)
	// 10 byte origin sender
	proto.suffixOfBlank(buf, response.OrigSender, 10)
	// 8 byte message length
	proto.prefixOfZero(buf, response.Payload.Len(), 8)
	// 8 byte control bits.
	buf.WriteString(response.CtrlBits)
	// 4 byte AreaCode
	proto.suffixOfBlank(buf, response.AreaCode, 4)
	// 4 byte fixed version
	buf.WriteString("0001")
	// 20 byte
	proto.suffixOfBlank(buf, response.MessageID, 20)
	// 20 byte
	proto.suffixOfBlank(buf, response.MessageRefID, 20)
	// 45 byte
	proto.suffixOfBlank(buf, response.Reserve, 45)
	// 1 byte end flag
	buf.WriteString("}")

	if response.Payload.Len() > 0 {
		// 2. write payload bytes
		buf.Write(response.Payload.Bytes())
	}

	return encrypt(ctx, buf, packetLen)
}

// prefixOfZero Appends '0' character until 10 bytes are satisfied
// eg: 0000000064, length 64, append prefix 00000000
func (proto *Protocol) prefixOfZero(buf buffer.IoBuffer, num int, max int) {
	rayLen := strconv.Itoa(num)
	if count := max - len(rayLen); count > 0 {
		for i := 0; i < count; i++ {
			buf.WriteString("0")
		}
	}
	buf.WriteString(rayLen)
}

// suffixOfBlank Appends ' ' character until max length are satisfied
func (proto *Protocol) suffixOfBlank(buf buffer.IoBuffer, val string, max int) {
	buf.WriteString(val)
	if count := max - len(val); count > 0 {
		for i := 0; i < count; i++ {
			buf.WriteString(" ")
		}
	}
}

func resolveHeaders(data []byte, h *common.Header) error {
	if len(data) <= 0 {
		return nil
	}

	v, err := parseXmlHeader(data)
	if err != nil {
		log.DefaultLogger.Errorf("failed to resolve beis proto header, err %v, data: %s", err, string(data))
		return err
	}

	var (
		code  string
		scene string
		flag  = requestFlag
	)

	if len(*v) <= 0 {
		// should never happen
		log.DefaultLogger.Warnf("resolved empty beis proto header, data: %s", string(data))
	}

	// only response contains RetStatus
	if s, ok := (*v)[retStatusKey]; ok && s != "" {
		flag = responseFlag
	}

	code = (*v)[serviceCodeKey]
	scene = (*v)[serviceSceneKey]
	// inject service id: code + scene
	h.Set(serviceKey, code+scene)

	// resolve Document xmlns
	xmlns := ""
	ds := string(data)
	i := strings.Index(ds, "xmlns=\"")
	if i >= 0 {
		xmlns = string(data[i+7 : i+7+20 /** `****.****.****.**.**` **/])
		// hijack response used.
		h.Set(xmlnsKey, xmlns)
	}

	if s, _ := h.Get(serviceKey); s == "" {
		// resolve xml namespace
		s = strings.ReplaceAll(xmlns, ".", "")
		h.Set(serviceKey, strings.ToUpper(s))
	}

	// inject other head if required.
	h.Set(requestTypeKey, flag)

	// update header unchanged, avoid encode again.
	h.Changed = false

	return nil
}

// parseXmlHeader decode xml header
func parseXmlHeader(data []byte) (*SystemHeader, error) {
	xmlBody := string(data)
	index := strings.Index(xmlBody, startHeader)
	header := &SystemHeader{}
	// parse header key value
	if index >= 0 {
		headerEndIndex := strings.Index(xmlBody, endHeader)
		xmlHeader := xmlBody[index : headerEndIndex+len(endHeader)]
		if xmlHeader != "" {
			xmlDecoder := xml.NewDecoder(bytes.NewBufferString(xmlHeader))
			if err := xmlDecoder.Decode(&header); err != nil {
				return nil, err
			}
		}
	}
	return header, nil
}
