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
	"mosn.io/api"
)

var Xr api.ProtocolName = "xr" // protocol

const (
	startHeader          = "<Header>"
	endHeader            = "</Header>"
	channelIdKey         = "ChannelId"
	requestIdKey         = "SteamId"
	externalReferenceKey = "ExternalReference"
	serviceCodeKey       = "ServiceCode"

	// request or response
	requestTypeKey = "RequestType"
	requestFlag    = "0" // 0 request
	responseFlag   = "1" // 1 response

	ResponseStatusSuccess uint32 = 0 // 0 response status
	RequestHeaderLen      int    = 8 // fix 8 byte header length
)

// StreamId query mapping stream id
func (proto *Proto) StreamId(ctx context.Context, key string) (val uint64, found bool) {
	val, found = proto.streams.Get(key)
	return
}

// PutStreamId put mapping stream id
func (proto *Proto) PutStreamId(ctx context.Context, key string, val uint64) (err error) {
	err = proto.streams.Put(key, val)
	return err
}

func (proto *Proto) RemoveStreamId(ctx context.Context, key string) {
	proto.streams.Remove(key)
	return
}
