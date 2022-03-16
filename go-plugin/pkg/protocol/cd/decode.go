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
	"context"
	"errors"
	"fmt"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/pkg/buffer"
	"strconv"
	"strings"
)

func (proto *Protocol) decodeRequest(ctx context.Context, buf api.IoBuffer, header *common.Header) (interface{}, error) {
	data := buf.Bytes()

	rawLen := strings.TrimLeft(string(data[0:10]), "0")

	var err error
	var packetLen = 0
	if rawLen != "" {
		// resolve fix 10 byte length
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to decode cd proto request len %d, err: %v", packetLen, err))
		}
	}

	totalLen := 10 /** fixed 10 byte len */ + packetLen
	// Read the complete packet data from the connection
	buf.Drain(totalLen)

	request := &Request{}

	// decode request field
	request.Header = *header

	request.Timeout = defaultTimeout

	request.Data = buffer.GetIoBuffer(totalLen)
	request.Data.Write(data[0:totalLen])

	payload := request.Data.Bytes()[10:totalLen]
	request.Payload = buffer.NewIoBufferBytes(payload)

	return request, nil
}

func (proto *Protocol) decodeResponse(ctx context.Context, buf api.IoBuffer, header *common.Header) (interface{}, error) {
	data := buf.Bytes()

	rawLen := strings.TrimLeft(string(data[0:10]), "0")

	var err error
	var packetLen = 0
	if rawLen != "" {
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to decode cd protoco request len %d, err: %v", packetLen, err))
		}
	}

	totalLen := 10 /** fixed 10 byte len */ + packetLen
	// Read the complete packet data from the connection
	buf.Drain(totalLen)

	response := &Response{}

	// decode request field
	response.Header = *header

	// Check whether the packet succeeds based on the actual packet
	// response.Status = 0

	response.Data = buffer.GetIoBuffer(totalLen)
	response.Data.Write(data[:totalLen])

	payload := response.Data.Bytes()[10:totalLen]
	response.Payload = buffer.NewIoBufferBytes(payload)

	return response, nil
}
