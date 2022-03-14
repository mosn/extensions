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
	"bytes"
	"context"
	"errors"
	"fmt"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/extensions/go-plugin/pkg/common/safe"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
)

// BeisProtocol protocol format: 128 byte header + string body
// <Document>
//    <SysHead>
//        <key>...</key>
//        <RetStatus>1</RetStatus>
//        <Ret>
//           <RetMsg>...</RetMsg>
//           <RetCode>..</RetCode>
//        </Ret>
//    </SysHead>
//    <AppHead>
//        <key>...</key>
//    </AppHead>
//    <key>...</key>
//    <details>
//        <xx>
//             <key>...</key>
//        </xx>
//    </details>
//  </Document>

type Protocol struct {
	streams safe.IntMap
}

func (proto *Protocol) Name() api.ProtocolName {
	return Beis
}

func (proto *Protocol) Encode(ctx context.Context, model interface{}) (api.IoBuffer, error) {
	switch frame := model.(type) {
	case *Request:
		return proto.encodeRequest(ctx, frame)
	case *Response:
		return proto.encodeResponse(ctx, frame)
	default:
		log.DefaultLogger.Errorf("[protocol][beis] encode with unknown command : %+v", model)
		return nil, errors.New("unknown command type")
	}
}

func (proto *Protocol) Decode(ctx context.Context, buf api.IoBuffer) (interface{}, error) {

	bLen := buf.Len()
	data := buf.Bytes()

	if bLen < RequestHeaderLen /** beis header length*/ {
		return nil, nil
	}

	var packetLen = 0
	var err error

	rawLen := strings.TrimLeft(string(data[MessageLengthIndex:MessageLengthIndex+8]), "0")
	if rawLen != "" {
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to decode beis proto package len %d, err: %v", packetLen, err))
		}
	}

	totalLen := RequestHeaderLen /** fixed 128 byte header len */ + packetLen
	// expected full message length
	if bLen < totalLen {
		return nil, nil
	}

	// decode buf xml body if encrypt
	ioBuf, err := decrypt(ctx, buf, totalLen)
	if err != nil {
		return nil, err
	}

	rpcHeader := common.Header{}
	resolveHeaders(ioBuf.Bytes()[RequestHeaderLen:RequestHeaderLen+totalLen], &rpcHeader)

	frameType, _ := rpcHeader.Get(requestTypeKey)
	switch frameType {
	case requestFlag:
		return proto.decodeRequest(ctx, ioBuf, &rpcHeader)
	case responseFlag:
		return proto.decodeResponse(ctx, ioBuf, &rpcHeader)
	default:
		return nil, fmt.Errorf("decode beis rpc Error, unknown request type = %s", frameType)
	}
}

// decrypt the body and return a new BUF if encryption exists
func decrypt(ctx context.Context, buf api.IoBuffer, totalLen int) (api.IoBuffer, error) {
	data := buf.Bytes()
	ctrlBits := data[26:34] // ctrl_bits
	cryptType := string(ctrlBits[0])

	switch cryptType {
	case "0": // plain text
		return buf, nil
	case "1": // xor ??
		buf.Drain(totalLen)
		// todo need to be implement
		// 1. decode xml body
		// 2. replace ioBuf MessageLength field, offset: [18:18+8], need append prefix '0'
		return nil, errors.New("unimplemented encrypt algorithm")
	case "2": // 3DS
		buf.Drain(totalLen)
		// todo need to be implement
		// 1. decode xml body
		// 2. replace ioBuf MessageLength field, offset: [18:18+8], need append prefix '0'
		return nil, errors.New("unimplemented encrypt 3DS algorithm")
	default:
		return nil, errors.New("unknown encrypt algorithm")
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
	return api.PingPong
}

func (proto *Protocol) EnableWorkerPool() bool {
	return false
}

func (proto *Protocol) GenerateRequestID(streamID *uint64) uint64 {
	return 0
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
	resolveHeaders(buf.Bytes()[10:10+len(body)], &rpcHeader)

	resp := NewRpcResponse(&rpcHeader, buf)
	return resp
}
