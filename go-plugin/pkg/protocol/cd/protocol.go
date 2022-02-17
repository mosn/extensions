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
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mosn/extensions/go-plugin/pkg/common"
	"github.com/mosn/extensions/go-plugin/pkg/common/safe"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
	"sync/atomic"
)

// CdProtocol protocol format: 10 byte length + string body
// <service>
//    <sys-header>
//        <data name="SYS_HEAD">
//            <struct>
//                <data name="field_name">
//                     <field length=int, type=string>...</field>
//                </data>
//                <data name="field_name_array">
//                     <array>
//                     	   <struct>
//                              <data name="field_name">
//                                 <field length=int, type=string>...</field>
//                              </data>
//                     	   </struct>
//                     </array>
//                </data>
//            </struct>
//        </data>
//    </sys-header>
//    <app-header>
//        <data name="APP_HEAD">
//            <struct>
//                <data name="field_name">
//                     <field length=int, type=string>...</field>
//                </data>
//                <data name="field_name_array">
//                     <array>
//                     	   <struct>
//                              <data name="field_name">
//                                 <field length=int, type=string>...</field>
//                              </data>
//                     	   </struct>
//                     </array>
//                </data>
//            </struct>
//        </data>
//    </app-header>
//    <local-header>
//        <data name="APP_HEAD">
//            <struct />
//        </data>
//    </local-header>
//    <Body>
//        <data name="field_name">
//             <field length=int, type=string>...</field>
//        </data>
//        <data name="field_name_array"> optional
//            <array>
//               <struct>
//                    <data name="field_name">
//                       <field length=int, type=string>...</field>
//                    </data>
//               </struct>
//            </array>
//        </data>
//        <data name="field_name_array"> optional
//            <struct>
//                <data name="field_name">
//                   <field length=int, type=string>...</field>
//                </data>
//            </struct>
//        </data>
//    </Body>
//  </service>
//
// ------------------ request example ---------------------------
// EXT_REF: Business requests are replaced automatically
//

type Protocol struct {
	streams safe.IntMap
}

func (proto *Protocol) Name() api.ProtocolName {
	return Cd
}

func (proto *Protocol) Encode(ctx context.Context, model interface{}) (api.IoBuffer, error) {
	switch frame := model.(type) {
	case *Request:
		return proto.encodeRequest(ctx, frame)
	case *Response:
		return proto.encodeResponse(ctx, frame)
	default:
		log.DefaultLogger.Errorf("[protocol][cd] encode with unknown command : %+v", model)
		return nil, errors.New("unknown command type")
	}
}

func (proto *Protocol) Decode(ctx context.Context, buf api.IoBuffer) (interface{}, error) {

	bLen := buf.Len()
	data := buf.Bytes()

	if bLen < 10 /** cd header length*/ {
		return nil, nil
	}

	var packetLen = 0
	var err error

	rawLen := strings.TrimLeft(string(data[0:10]), "0")
	if rawLen != "" {
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to decode cd ptoco package len %d, err: %v", packetLen, err))
		}
	}

	// expected full message length
	if bLen < packetLen {
		return nil, nil
	}

	totalLen := 10 /** fixed 10 byte len */ + packetLen

	rpcHeader := common.Header{}
	injectHeaders(data[10:totalLen], &rpcHeader)

	frameType, _ := rpcHeader.Get(requestTypeKey)
	switch frameType {
	case requestFlag:
		return proto.decodeRequest(ctx, buf, &rpcHeader)
	case responseFlag:
		return proto.decodeResponse(ctx, buf, &rpcHeader)
	default:
		return nil, fmt.Errorf("decode cd rpc Error, unkownen request type = %s", frameType)
	}
}

// Trigger heartbeat detect.
func (proto *Protocol) Trigger(context context.Context, requestId uint64) api.XFrame {
	return nil
}

func (proto *Protocol) Reply(context context.Context, request api.XFrame) api.XRespFrame {
	return nil
}

// Hijack hijack request, maybe timeout
func (proto *Protocol) Hijack(context context.Context, request api.XFrame, statusCode uint32) api.XRespFrame {
	resp := proto.hijackResponse(request, statusCode)

	return resp

}

func (proto *Protocol) Mapping(httpStatusCode uint32) uint32 {
	return httpStatusCode
}

// PoolMode returns whether ping-pong or multiplex
func (proto *Protocol) PoolMode() api.PoolMode {
	return api.Multiplex
}

func (proto *Protocol) EnableWorkerPool() bool {
	return true
}

func (proto *Protocol) GenerateRequestID(streamID *uint64) uint64 {
	return atomic.AddUint64(streamID, 1)
}

// hijackResponse build hijack response
func (proto *Protocol) hijackResponse(request api.XFrame, statusCode uint32) *Response {
	req := request.(*Request)
	body := req.Payload.String()

	var bodyBuf bytes.Buffer
	headerIndex := strings.Index(body, startHeader)
	if headerIndex >= 0 {
		bodyBuf.WriteString(body[:headerIndex+len(startHeader)])
		bodyBuf.WriteString("<Response>")
		bodyBuf.WriteString("<ReturnCode>")
		bodyBuf.WriteString(strconv.Itoa(int(statusCode)))
		bodyBuf.WriteString("</ReturnCode>")
		bodyBuf.WriteString("<ReturnMessage>此请求被劫持，code: ")
		bodyBuf.WriteString(strconv.Itoa(int(statusCode)))
		bodyBuf.WriteString("</ReturnMessage>")
		bodyBuf.WriteString("</Response>")
		bodyBuf.WriteString(body[headerIndex+len(startHeader):])
	}
	body = bodyBuf.String()
	// replace request type -> response
	body = strings.ReplaceAll(body, "<RequestType>0</RequestType>", "<RequestType>1</RequestType>")

	// 10 byte length + string body
	buf := buffer.GetIoBuffer(10 + len(body))
	proto.prefixOfZero(buf, len(body))
	buf.WriteString(body)

	// response header
	rpcHeader := common.Header{}
	injectHeaders(buf.Bytes()[10:10+len(body)], &rpcHeader)

	resp := NewRpcResponse(&rpcHeader, buf)
	return resp
}
