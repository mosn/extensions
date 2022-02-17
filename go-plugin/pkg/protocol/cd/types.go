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
	"encoding/xml"
	"mosn.io/api"
)

var Cd api.ProtocolName = "cd" // protocol

const (
	startHeader = "<sys-header>"
	endHeader   = "</sys-header>"

	requestIdKey = "RequestId"

	serviceCodeKey = "ServiceCode"

	// request or response
	requestTypeKey = "RequestType"
	requestFlag    = "0" // 0 request
	responseFlag   = "1" // 1 response

	ResponseStatusSuccess uint32 = 0  // 0 response status
	RequestHeaderLen      int    = 10 // fix 10 byte header length
)

// StreamId query mapping stream id
func (proto *Protocol) StreamId(ctx context.Context, key string) (val uint64, found bool) {
	val, found = proto.streams.Get(key)
	return
}

// PutStreamId put mapping stream id
func (proto *Protocol) PutStreamId(ctx context.Context, key string, val uint64) (err error) {
	err = proto.streams.Put(key, val)
	return err
}

func (proto *Protocol) RemoveStreamId(ctx context.Context, key string) {
	proto.streams.Remove(key)
	return
}

// SystemHeader cd protocol sys-header
// <service>
//    <sys-header>
type SystemHeader struct {
	XMLName  xml.Name   `xml:"data"`
	Name     string     `xml:"name,attr"`
	WrapData []WrapData `xml:"struct>data"`
}

// WrapData cd protocol data>struct
// <data name="SYS_HEAD">
//    <struct>
type WrapData struct {
	XMLName    xml.Name    `xml:"data"`
	Name       string      `xml:"name,attr"`
	Field      *Field      `xml:"field,omitempty"`
	ArrayField *[]WrapData `xml:"array>struct>data,omitempty"`
}

// Field cd protocol filed
// <data>
//    <field
type Field struct {
	Length int    `xml:"length,attr"`
	Scale  int    `xml:"scale,attr"`
	Type   string `xml:"type,attr"`
	Value  string `xml:",chardata"`
}
