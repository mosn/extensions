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
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
)

const defaultTimeout = 30000 // default request timeout(30 seconds).

type Request struct {
	common.Header              // request key value pair
	BeginFlag     string       // 8 byte, 起始标识
	OrigSender    string       // 10 byte, 发起系统标识
	Length        uint64       // 8 byte, 除报文头外的报文数据长度
	CtrlBits      string       // 8 byte
	AreaCode      string       // 4 byte
	VersionID     string       // 4 byte, fixed: 0001
	MessageID     string       // 20 byte
	MessageRefID  string       // 20 byte
	Reserve       string       // 45 byte
	EndFlag       string       // 1 byte, 结束标识
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
	// Request not implement request id so that request and response must be ping pong
	return 0
}

func (r *Request) SetRequestId(id uint64) {
	// Request not implement request id so that request and response must be ping pong
}

// check command implement api interface.
var _ api.XFrame = &Request{}
var _ api.XRespFrame = &Response{}

type Response struct {
	common.Header              // response key value pair
	BeginFlag     string       // 8 byte, 起始标识
	OrigSender    string       // 10 byte, 发起系统标识
	Length        uint64       // 8 byte, 除报文头外的报文数据长度
	CtrlBits      string       // 8 byte
	AreaCode      string       // 4 byte
	VersionID     string       // 4 byte, fixed: 0001
	MessageID     string       // 20 byte
	MessageRefID  string       // 20 byte
	Reserve       string       // 45 byte
	EndFlag       string       // 1 byte, 结束标识
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
	// Response not implement request id so that request and response must be ping pong
	return 0
}

func (r *Response) SetRequestId(id uint64) {
	// Response not implement request id so that request and response must be ping pong
}

func (r *Response) GetStatusCode() uint32 {
	// Response not implement status code
	return 0
}

func hash(s string) uint64 {
	var h uint64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 16777619
	}
	return h
}
