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
	"context"
	"errors"
	"fmt"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/pkg/buffer"
	"strconv"
	"strings"
)

func (proto *XrProtocol) decodeRequest(ctx context.Context, buf api.IoBuffer, header *common.Header) (interface{}, error) {
	bufLen := buf.Len()
	data := buf.Bytes()

	rawLen := strings.TrimLeft(string(data[0:8]), "0")

	var err error
	var packetLen = 0
	if rawLen != "" {
		// resolve fix 8 byte length
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to decode request len %d, err: %v", packetLen, err))
		}
	}

	// The buf does not contain the complete packet length,
	// So we wait for the next decoder notification.
	if bufLen < packetLen {
		return nil, nil
	}

	totalLen := 8 /** fixed 8 byte len */ + packetLen
	// Read the complete packet data from the connection
	buf.Drain(totalLen)

	// rpcBufCtx := bufferByContext(ctx)

	request := &Request{}

	// decode request field
	request.Header = *header
	request.RequestId, _ = header.Get(requestIdKey)
	if val, found := proto.StreamId(ctx, request.RequestId); found {
		request.SteamId = val
	}
	request.Timeout = defaultTimeout

	request.Data = buffer.GetIoBuffer(totalLen)
	request.Data.Write(data[:totalLen])

	payload := request.Data.Bytes()[8:totalLen]
	request.Payload = buffer.NewIoBufferBytes(payload)

	return request, nil
}

func (proto *XrProtocol) decodeResponse(ctx context.Context, buf api.IoBuffer, header *common.Header) (interface{}, error) {
	bufLen := buf.Len()
	data := buf.Bytes()

	rawLen := strings.TrimLeft(string(data[0:8]), "0")

	var err error
	var packetLen = 0
	if rawLen != "" {
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to decode request len %d, err: %v", packetLen, err))
		}
	}

	// The buf does not contain the complete packet length,
	// So we wait for the next decoder notification.
	if bufLen < packetLen {
		return nil, nil
	}

	totalLen := 8 /** fixed 8 byte len */ + packetLen
	// Read the complete packet data from the connection
	buf.Drain(totalLen)

	// rpcBufCtx := bufferByContext(ctx)

	response := &Response{}

	// decode request field
	response.Header = *header
	response.RequestId, _ = header.Get(requestIdKey)
	if val, found := proto.StreamId(ctx, response.RequestId); found {
		response.SteamId = val
		// remove stream id, help gc
		proto.RemoveStreamId(ctx, response.RequestId)
	}

	// Check whether the packet succeeds based on the actual packet
	// response.Status = 0

	response.Data = buffer.GetIoBuffer(totalLen)
	response.Data.Write(data[:totalLen])

	payload := response.Data.Bytes()[8:totalLen]
	response.Payload = buffer.NewIoBufferBytes(payload)

	return response, nil
}
