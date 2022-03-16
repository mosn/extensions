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
	"context"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/pkg/buffer"
	"strconv"
	"strings"
)

func (proto *Protocol) decodeRequest(ctx context.Context, buf api.IoBuffer, header *common.Header) (interface{}, error) {
	data := buf.Bytes()

	rawLen := strings.TrimLeft(string(data[MessageLengthIndex:MessageLengthIndex+8]), "0")
	packetLen, _ := strconv.Atoi(rawLen)

	totalLen := RequestHeaderLen /** fixed 128 byte header len */ + packetLen

	// Read the complete packet data from the connection
	buf.Drain(totalLen)

	request := &Request{
		BeginFlag:    string(data[0:8]),
		OrigSender:   string(data[8:18]),
		Length:       uint64(packetLen),
		CtrlBits:     string(data[26:34]),
		AreaCode:     string(data[34:38]),
		VersionID:    string(data[38:42]),
		MessageID:    string(data[42:62]),
		MessageRefID: string(data[62:82]),
		Reserve:      string(data[82:127]),
		EndFlag:      string(data[127:128]),
	}

	// decode request field
	request.Header = *header
	request.Timeout = defaultTimeout

	request.Data = buffer.GetIoBuffer(totalLen)
	request.Data.Write(data[:totalLen])

	payload := request.Data.Bytes()[RequestHeaderLen:totalLen]
	request.Payload = buffer.NewIoBufferBytes(payload)

	return request, nil
}

func (proto *Protocol) decodeResponse(ctx context.Context, buf api.IoBuffer, header *common.Header) (interface{}, error) {
	data := buf.Bytes()

	rawLen := strings.TrimLeft(string(data[MessageLengthIndex:MessageLengthIndex+8]), "0")
	packetLen, _ := strconv.Atoi(rawLen)

	totalLen := RequestHeaderLen /** fixed 128 byte header len */ + packetLen

	// Read the complete packet data from the connection
	buf.Drain(totalLen)

	response := &Response{
		BeginFlag:    string(data[0:8]),
		OrigSender:   string(data[8:18]),
		Length:       uint64(packetLen),
		CtrlBits:     string(data[26:34]),
		AreaCode:     string(data[34:38]),
		VersionID:    string(data[38:42]),
		MessageID:    string(data[42:62]),
		MessageRefID: string(data[62:82]),
		Reserve:      string(data[82:127]),
		EndFlag:      string(data[127:128]),
	}

	// decode request field
	response.Header = *header

	response.Data = buffer.GetIoBuffer(totalLen)
	response.Data.Write(data[:totalLen])

	payload := response.Data.Bytes()[RequestHeaderLen:totalLen]
	response.Payload = buffer.NewIoBufferBytes(payload)

	return response, nil
}
