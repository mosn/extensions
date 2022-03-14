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
)

var proto = &Protocol{}

// NewRpcRequest is a utility function which build rpc Request object of xr protocol.
func NewRpcRequest(headers *common.Header, data api.IoBuffer) *Request {
	frame, err := proto.decodeRequest(context.Background(), data, headers)
	if err != nil {
		return nil
	}
	request, ok := frame.(*Request)
	if !ok {
		return nil
	}

	// set headers
	if headers != nil {
		headers.Range(func(key, value string) bool {
			request.Header.Set(key, value)
			return true
		})
	}
	return request
}

// NewRpcResponse is a utility function which build rpc Response object of xr protocol.
func NewRpcResponse(headers *common.Header, data api.IoBuffer) *Response {
	frame, err := proto.decodeResponse(context.Background(), data, headers)
	if err != nil {
		return nil
	}
	response, ok := frame.(*Response)
	if !ok {
		return nil
	}

	// set headers
	if headers != nil {
		headers.Range(func(key, value string) bool {
			response.Header.Set(key, value)
			return true
		})
	}
	return response
}
