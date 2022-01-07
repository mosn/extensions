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

package bolt

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"mosn.io/api"

	"mosn.io/pkg/log"
)

/**
 * Request command protocol for v1
 * 0     1     2           4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |proto| type| cmdcode   |ver2 |   requestID           |codec|        timeout        |  classLen |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 * |headerLen  | contentLen            |                             ... ...                       |
 * +-----------+-----------+-----------+                                                                                               +
 * |               className + header  + content  bytes                                            |
 * +                                                                                               +
 * |                               ... ...                                                         |
 * +-----------------------------------------------------------------------------------------------+
 *
 * proto: code for protocol
 * type: request/response/request oneway
 * cmdcode: code for remoting command
 * ver2:version for remoting command
 * requestID: id of request
 * codec: code for codec
 * headerLen: length of header
 * contentLen: length of content
 *
 * Response command protocol for v1
 * 0     1     2     3     4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |proto| type| cmdcode   |ver2 |   requestID           |codec|respstatus |  classLen |headerLen  |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 * | contentLen            |                  ... ...                                              |
 * +-----------------------+                                                                       +
 * |                         className + header  + content  bytes                                  |
 * +                                                                                               +
 * |                               ... ...                                                         |
 * +-----------------------------------------------------------------------------------------------+
 * respstatus: response status
 */

type BoltProtocol struct{}

// types.Protocol
func (proto BoltProtocol) Name() api.ProtocolName {
	return ProtocolName
}

func (proto BoltProtocol) Encode(ctx context.Context, model interface{}) (api.IoBuffer, error) {
	switch frame := model.(type) {
	case *Request:
		return encodeRequest(ctx, frame)
	case *Response:
		return encodeResponse(ctx, frame)
	default:
		log.DefaultLogger.Errorf("[protocol][bolt] encode with unknown command : %+v", model)
		return nil, api.ErrUnknownType
	}
}

func (proto BoltProtocol) Decode(ctx context.Context, data api.IoBuffer) (interface{}, error) {
	if data.Len() >= LessLen {
		cmdType := data.Bytes()[1]

		switch cmdType {
		case CmdTypeRequest:
			return decodeRequest(ctx, data, false)
		case CmdTypeRequestOneway:
			return decodeRequest(ctx, data, true)
		case CmdTypeResponse:
			return decodeResponse(ctx, data)
		default:
			// unknown cmd type
			return nil, fmt.Errorf("Decode Error, type = %s, value = %d", UnKnownCmdType, cmdType)
		}
	}

	return nil, nil
}

// Heartbeater
func (proto BoltProtocol) Trigger(ctx context.Context, requestId uint64) api.XFrame {
	return &Request{
		RequestHeader: RequestHeader{
			Protocol:  ProtocolCode,
			CmdType:   CmdTypeRequest,
			CmdCode:   CmdCodeHeartbeat,
			Version:   1,
			RequestId: uint32(requestId),
			Codec:     Hessian2Serialize,
			Timeout:   -1,
		},
	}
}

func (proto BoltProtocol) Reply(ctx context.Context, request api.XFrame) api.XRespFrame {
	return &Response{
		ResponseHeader: ResponseHeader{
			Protocol:       ProtocolCode,
			CmdType:        CmdTypeResponse,
			CmdCode:        CmdCodeHeartbeat,
			Version:        ProtocolVersion,
			RequestId:      uint32(request.GetRequestId()),
			Codec:          Hessian2Serialize,
			ResponseStatus: ResponseStatusSuccess,
		},
	}
}

// Hijacker
func (proto BoltProtocol) Hijack(ctx context.Context, request api.XFrame, statusCode uint32) api.XRespFrame {
	return &Response{
		ResponseHeader: ResponseHeader{
			Protocol:       ProtocolCode,
			CmdType:        CmdTypeResponse,
			CmdCode:        CmdCodeRpcResponse,
			Version:        ProtocolVersion,
			RequestId:      0,                 // this would be overwrite by stream layer
			Codec:          Hessian2Serialize, //todo: read default codec from config
			ResponseStatus: uint16(statusCode),
		},
	}
}

func (proto BoltProtocol) Mapping(httpStatusCode uint32) uint32 {
	switch httpStatusCode {
	case http.StatusOK:
		return uint32(ResponseStatusSuccess)
	case api.RouterUnavailableCode:
		return uint32(ResponseStatusNoProcessor)
	case api.NoHealthUpstreamCode:
		return uint32(ResponseStatusConnectionClosed)
	case api.UpstreamOverFlowCode:
		return uint32(ResponseStatusServerThreadpoolBusy)
	case api.CodecExceptionCode:
		//Decode or Encode Error
		return uint32(ResponseStatusCodecException)
	case api.DeserialExceptionCode:
		//Hessian Exception
		return uint32(ResponseStatusServerDeserialException)
	case api.TimeoutExceptionCode:
		//Response Timeout
		return uint32(ResponseStatusTimeout)
	case api.PermissionDeniedCode:
		//Response Permission Denied
		// bolt protocol do not have a permission deny code, use server exception
		return uint32(ResponseStatusServerException)
	default:
		return uint32(ResponseStatusUnknown)
	}
}

// PoolMode returns whether pingpong or multiplex
func (proto BoltProtocol) PoolMode() api.PoolMode {
	return api.Multiplex
}

func (proto BoltProtocol) EnableWorkerPool() bool {
	return true
}

func (proto BoltProtocol) GenerateRequestID(streamID *uint64) uint64 {
	return atomic.AddUint64(streamID, 1)
}
