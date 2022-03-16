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
	"io"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
)

func (proto *XrProtocol) encodeRequest(ctx context.Context, request *Request) (api.IoBuffer, error) {

	packetLen := 8 /** fixed 8 byte length */ + request.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 8 byte length + body
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
			log.DefaultLogger.Debugf("xr proto mapping streamId: %d -> %d", request.RequestId, request.SteamId.(uint64))
		}
	}

	return buf, nil
}

func (proto *XrProtocol) encodeResponse(ctx context.Context, response *Response) (api.IoBuffer, error) {

	packetLen := 8 /** fixed 8 byte length */ + response.Payload.Len()
	buf := buffer.GetIoBuffer(packetLen)

	// 1. write 8 byte length + body
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

// prefixOfZero Appends '0' character until 8 bytes are satisfied
// eg: 000064, length 64, append prefix 0000
func (proto *XrProtocol) prefixOfZero(buf buffer.IoBuffer, payloadLen int) {
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

	header, err := parseXmlHeader(data)
	if err != nil {
		return err
	}

	channelId := header[channelIdKey]
	extRef := header[externalReferenceKey]
	code := header[serviceCodeKey]
	flag := header[requestTypeKey]

	// inject request id: channelId + extRef
	h.Set(requestIdKey, channelId+extRef)

	// inject service id
	h.Set(serviceCodeKey, code)

	// inject other head if required.
	h.Set(requestTypeKey, flag)

	// update header unchanged, avoid encode again.
	h.Changed = false

	return nil
}

// parseXmlHeader decode xml header
func parseXmlHeader(data []byte) (XmlHeader, error) {
	xmlBody := string(data)
	index := strings.Index(xmlBody, startHeader)
	header := XmlHeader{}
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

// XmlHeader xml key value pair.
// Protocol-specific, depending on
// traditional protocol data structures
type XmlHeader map[string]string

type KeyValueEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XmlHeader) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(KeyValueEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *XmlHeader) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = XmlHeader{}
	for {
		var e KeyValueEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}

	return nil
}
