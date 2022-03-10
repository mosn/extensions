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
		if request.RequestId == "" {
			// try query business id from payload.
			payload := request.Payload.Bytes()
			request.RequestId = fetchId(payload)
		}

		// 2. write payload bytes
		buf.Write(request.Payload.Bytes())
	}

	// If sidecar replaces the ID, we associate the ID with the business ID
	// When the response is received, streamId is restored correctly.
	if request.SteamId != nil {
		proto.PutStreamId(ctx, request.RequestId, request.SteamId.(uint64))
		// record debug mapping stream info.
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("cd proto mapping streamId: %d -> %d", request.RequestId, request.SteamId.(uint64))
		}
	}

	return buf, nil
}

func (proto *Protocol) encodeResponse(ctx context.Context, response *Response) (api.IoBuffer, error) {

	packetLen := 10 /** fixed 10 byte length */ + response.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 10 byte length + body
	proto.prefixOfZero(buf, response.Payload.Len())
	if response.Payload.Len() > 0 {
		if response.RequestId == "" {
			payload := response.Payload.Bytes()
			response.RequestId = fetchId(payload)
		}

		// 2. write payload bytes
		buf.Write(response.Payload.Bytes())
	}

	// remove associate the ID with the business ID if exists.
	if _, found := proto.StreamId(ctx, response.RequestId); found {
		// remove stream id, help gc
		proto.RemoveStreamId(ctx, response.RequestId)
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

	v, err := parseXmlHeader(data)
	if err != nil {
		log.DefaultLogger.Errorf("failed to resolve cd proto header, err %v, data: %s", err, string(data))
		return err
	}

	var (
		code  string
		scene string
		reqId string
		flag  = requestFlag
	)

	if len(v.WrapData) > 0 {
		for _, d := range v.WrapData {
			if d.Field != nil { // plain field
				switch d.Name {
				case "SERVICE_CODE":
					code = d.Field.Value
				case "SERVICE_SCENE":
					scene = d.Field.Value
				case "SERVICE_REQUEST_ID": // todo need to be modified.
					reqId = d.Field.Value
				}
			} else if d.ArrayField != nil {
				for _, f := range *d.ArrayField {
					if f.Field != nil {
						// struct field
						switch f.Name {
						case "RET":
							flag = responseFlag
						}
					}
				}
			}
		}
	} else {
		// should never happen
		log.DefaultLogger.Warnf("resolved empty cd proto header, data: %s", string(data))
	}

	// check request id must exist
	if reqId == "" {
		log.DefaultLogger.Warnf("cd proto header req id must exist, data: %s", string(data))
	}

	// inject request id
	h.Set(requestIdKey, reqId)

	// inject service id: code + scene
	h.Set(serviceCodeKey, code+scene)

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
