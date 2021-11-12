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

package xr

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/mosn/wasm-sdk/go-plugin/pkg/common"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/header"
	"strconv"
	"strings"
)

func (proto *Proto) encodeRequest(ctx context.Context, request *Request) (api.IoBuffer, error) {

	packetLen := 8 /** fixed 8 byte length */ + request.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 8 byte length + body
	proto.prefixOfZero(buf, request.Payload.Len())

	if request.Payload.Len() > 0 {
		if request.RequestId == "" {
			// try query business id from payload.
			payload := request.Payload.Bytes()
			request.RequestId = fetchId(payload[8 : 8+len(payload)])
		}

		// 2. write payload bytes
		buf.Write(request.Payload.Bytes())
	}

	// If sidecar replaces the ID, we associate the ID with the business ID
	// When the response is received, streamId is restored correctly.
	if request.SteamId != nil {
		proto.PutStreamId(ctx, request.RequestId, request.SteamId.(uint64))
	}

	return buf, nil
}

func (proto *Proto) encodeResponse(ctx context.Context, response *Response) (api.IoBuffer, error) {

	packetLen := 8 /** fixed 8 byte length */ + response.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 8 byte length + body
	proto.prefixOfZero(buf, response.Payload.Len())
	if response.Payload.Len() > 0 {
		if response.RequestId == "" {
			payload := response.Payload.Bytes()
			response.RequestId = fetchId(payload[8 : 8+len(payload)])
		}

		// 2. write payload bytes
		buf.Write(response.Payload.Bytes())
	}

	// If sidecar replaces the ID, we associate the ID with the business ID
	// When the response is received, streamId is restored correctly.
	if response.SteamId != nil {
		proto.PutStreamId(ctx, response.RequestId, response.SteamId.(uint64))
	}

	return buf, nil
}

// prefixOfZero Appends '0' character until 8 bytes are satisfied
// eg: 000064, length 64, append prefix 0000
func (proto *Proto) prefixOfZero(buf buffer.IoBuffer, payloadLen int) {
	rayLen := strconv.Itoa(payloadLen)
	if count := 8 - len(rayLen); count > 0 {
		for i := 0; i < count; i++ {
			buf.WriteString("0")
		}
	}
	buf.WriteString(rayLen)
}

func fetchId(data []byte) string {
	h := common.Header{}
	injectHeaders(data, &h)
	id, _ := h.Get(requestIdKey)
	return id
}

func injectHeaders(data []byte, h *common.Header) error {
	if len(data) <= 0 {
		return nil
	}

	xmlBody := string(data)
	index := strings.Index(xmlBody, startHeader)
	header := header.CommonHeader{}
	// parse header key value
	if index >= 0 {
		headerEndIndex := strings.Index(xmlBody, endHeader)
		xmlHeader := xmlBody[index : headerEndIndex+len(endHeader)]
		if xmlHeader != "" {
			xmlDecoder := xml.NewDecoder(bytes.NewBufferString(xmlHeader))
			if err := xmlDecoder.Decode(&header); err != nil {
				return err
			}
		}
	}

	channelId := header[channelIdKey]
	extRef := header[externalReferenceKey]
	code := header[serviceCodeKey]

	// inject request id: channelId + extRef
	h.Set(requestIdKey, channelId+extRef)

	// inject service id
	h.Set(serviceCodeKey, code)

	// inject other head if required.

	// update header unchanged, avoid encode again.
	h.Changed = false

	return nil
}
