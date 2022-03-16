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
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
)

const defaultTimeout = 10000 // default request timeout(10 seconds).

type Request struct {
	common.Header              // request key value pair
	RequestId     string       // request id (biz id)
	SteamId       interface{}  // sidecar request id (replaced by sidecar, uint64 or nil)
	Timeout       uint32       // request timeout
	Payload       api.IoBuffer // it refers to the service parameters of a packet
	Data          api.IoBuffer // full package bytes
	Changed       bool         // indicates whether the packet payload is modified
}

func (r *Request) IsHeartbeatFrame() bool {
	return false
}

func (r *Request) GetTimeout() int32 {
	return defaultTimeout
}

func (r *Request) GetHeader() api.HeaderMap {
	return r
}

func (r *Request) GetData() api.IoBuffer {
	return r.Payload
}

func (r *Request) SetData(data api.IoBuffer) {
	r.Payload = data
}

func (r *Request) GetStreamType() api.StreamType {
	return api.Request
}

func (r *Request) GetRequestId() uint64 {
	if r.SteamId != nil {
		return r.SteamId.(uint64)
	}

	// we don't care about it
	return hash(r.RequestId)
}

func (r *Request) SetRequestId(id uint64) {
	r.SteamId = id
}

// check command implement api interface.
var _ api.XFrame = &Request{}
var _ api.XRespFrame = &Response{}

type Response struct {
	common.Header              // response key value pair
	RequestId     string       // response id
	SteamId       interface{}  // sidecar request id (replaced by sidecar id)
	Status        uint32       // response status
	Data          api.IoBuffer // full package bytes
	Payload       api.IoBuffer // it refers to the service parameters of a packet
	Changed       bool         // indicates whether the packet payload is modified
}

func (r *Response) IsHeartbeatFrame() bool {
	return false
}

func (r *Response) GetTimeout() int32 {
	return defaultTimeout
}

func (r *Response) GetHeader() api.HeaderMap {
	return r
}

func (r *Response) GetData() api.IoBuffer {
	return r.Payload
}

func (r *Response) SetData(data api.IoBuffer) {
	r.Payload = data
}

func (r *Response) GetStreamType() api.StreamType {
	return api.Response
}

func (r *Response) GetRequestId() uint64 {
	if r.SteamId != nil {
		return r.SteamId.(uint64)
	}

	// we don't care about it
	return hash(r.RequestId)
}

func (r *Response) SetRequestId(id uint64) {
	r.SteamId = id
}

func (r *Response) GetStatusCode() uint32 {
	return r.Status
}

func hash(s string) uint64 {
	var h uint64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 16777619
	}
	return h
}
