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

package cd

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

	packetLen := 10 /** fixed 10 byte length */ + request.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 10 byte length + body
	proto.prefixOfZero(buf, request.Payload.Len())

	if request.Payload.Len() > 0 {
		// 2. write payload bytes
		buf.Write(request.Payload.Bytes())
	}

	return buf, nil
}

func (proto *Protocol) encodeResponse(ctx context.Context, response *Response) (api.IoBuffer, error) {

	packetLen := 10 /** fixed 10 byte length */ + response.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 10 byte length + body
	proto.prefixOfZero(buf, response.Payload.Len())
	if response.Payload.Len() > 0 {
		// 2. write payload bytes
		buf.Write(response.Payload.Bytes())
	}

	return buf, nil
}

// prefixOfZero Appends '0' character until 10 bytes are satisfied
// eg: 0000000064, length 64, append prefix 00000000
func (proto *Protocol) prefixOfZero(buf buffer.IoBuffer, payloadLen int) {
	rayLen := strconv.Itoa(payloadLen)
	if count := 10 - len(rayLen); count > 0 {
		for i := 0; i < count; i++ {
			buf.WriteString("0")
		}
	}
	buf.WriteString(rayLen)
}

func injectHeaders(data []byte, h *common.Header) error {
	if len(data) <= 0 {
		return nil
	}

	v, err := parseXmlHeader(data, startHeader, endHeader)
	if err != nil {
		log.DefaultLogger.Errorf("failed to resolve cd proto header, err %v, data: %s", err, string(data))
		return err
	}

	var (
		code          string
		scene         string
		consumerId    string
		consumerSeqNo string
		tranId        string
		branchId      string
		userId        string
		flag          = requestFlag
	)

	if len(v.WrapData) > 0 {
		for _, d := range v.WrapData {
			if d.Field != nil { // plain field
				switch d.Name {
				case serviceCodeKey:
					code = d.Field.Value
				case serviceSceneKey:
					scene = d.Field.Value
				case consumerIdKey:
					consumerId = d.Field.Value
				case consumerSeqNoKey:
					consumerSeqNo = d.Field.Value
				case tranIdKey:
					tranId = d.Field.Value
				case branchIdKey:
					branchId = d.Field.Value
				case userIdKey:
					userId = d.Field.Value
				case retStatusKey, retKey:
					flag = responseFlag
				}
			}
		}
	} else {
		// should never happen
		log.DefaultLogger.Warnf("resolved empty cd proto header, data: %s", string(data))
	}

	// inject service id: code + scene
	h.Set(serviceKey, code+scene)
	h.Set(serviceCodeKey, code)
	h.Set(serviceSceneKey, scene)

	// inject consumerIdKey
	h.Set(consumerIdKey, consumerId)

	// inject consumerSeqNoKey
	h.Set(consumerSeqNoKey, consumerSeqNo)

	// inject tranIdKey
	h.Set(tranIdKey, tranId)

	// inject branchIdKey
	h.Set(branchIdKey, branchId)

	// inject userIdKey
	h.Set(userIdKey, userId)

	// inject other head if required.
	h.Set(requestTypeKey, flag)

	// update header unchanged, avoid encode again.
	h.Changed = false

	return nil
}

// parseXmlHeader decode xml header
func parseXmlHeader(data []byte, startTag, endTag string) (*SystemHeader, error) {
	xmlBody := string(data)
	index := strings.Index(xmlBody, startTag)
	header := &SystemHeader{}
	// parse header key value
	if index >= 0 {
		headerEndIndex := strings.Index(xmlBody, endTag)
		xmlHeader := xmlBody[index+len(startTag) : headerEndIndex]
		if xmlHeader != "" {
			xmlDecoder := xml.NewDecoder(bytes.NewBufferString(xmlHeader))
			if err := xmlDecoder.Decode(&header); err != nil {
				return nil, err
			}
		}
	}
	return header, nil
}
