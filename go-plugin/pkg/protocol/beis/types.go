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
	"encoding/xml"
	"io"
	"mosn.io/api"
)

var Beis api.ProtocolName = "beis" // protocol

const (
	startHeader = "<SysHead>"
	endHeader   = "</SysHead>"

	serviceKey = "service"
	xmlnsKey   = "xmlns"

	serviceCodeKey  = "ServiceCode"
	serviceSceneKey = "ServiceScene"
	retStatusKey    = "RetStatus"

	// request or response
	requestTypeKey = "RequestType"
	requestFlag    = "0" // 0 request
	responseFlag   = "1" // 1 response

	beginFlag = "{BOBXML:"

	ResponseStatusSuccess uint32 = 0   // 0 response status
	RequestHeaderLen      int    = 128 // fix 128 byte header length
	MessageLengthIndex    int    = 18  // message content offset index
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

// SystemHeader decode sys-header key value pair
type SystemHeader map[string]string

type KeyValueEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m SystemHeader) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *SystemHeader) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = SystemHeader{}
	for {
		var e KeyValueEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}

	return nil
}
