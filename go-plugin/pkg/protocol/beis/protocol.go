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

	var bodyBuf bytes.Buffer

	// 1. write xml header
	bodyBuf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")

	// 2. write Document tag
	xmlns, _ := req.Get(xmlnsKey)
	bodyBuf.WriteString("<Document ")
	{
		// Document attribute
		bodyBuf.WriteString("xmlns=\"" + xmlns + "\">")
	}

	// 3. write SysHead
	bodyBuf.WriteString("<SysHead>")
	{
		bodyBuf.WriteString("<RetStatus>1</RetStatus>") // failed
		bodyBuf.WriteString("<Ret>")
		{
			code, message := mappingCode(statusCode)

			bodyBuf.WriteString("<RetCode>")
			bodyBuf.WriteString(code)
			bodyBuf.WriteString("</RetCode>")

			bodyBuf.WriteString("<RetMsg>")
			bodyBuf.WriteString(message)
			bodyBuf.WriteString("</RetMsg>")
		}
		bodyBuf.WriteString("</Ret>")
	}
	bodyBuf.WriteString("</SysHead>")

	// 4. write LOCAL_HEAD
	bodyBuf.WriteString("<LOCAL_HEAD/>")

	// 5. write AppHead
	bodyBuf.WriteString("<AppHead/>")
	bodyBuf.WriteString("</Document>")

	// 128 byte length + string body
	buf := buffer.GetIoBuffer(128 + bodyBuf.Len())

	// 1. write 128 byte length, 8 byte fixed begin flag
	buf.WriteString(beginFlag)
	// 10 byte origin sender
	proto.suffixOfBlank(buf, req.OrigSender, 10)
	// 8 byte message length
	proto.prefixOfZero(buf, bodyBuf.Len(), 8)
	// 8 byte control bits.
	buf.WriteString(req.CtrlBits)
	// 4 byte AreaCode
	proto.suffixOfBlank(buf, req.AreaCode, 4)
	// 4 byte fixed version
	buf.WriteString("0001")
	// 20 byte
	proto.suffixOfBlank(buf, req.MessageID, 20)
	// 20 byte
	proto.suffixOfBlank(buf, req.MessageRefID, 20)
	// 45 byte
	proto.suffixOfBlank(buf, req.Reserve, 45)
	// 1 byte end flag
	buf.WriteString("}")

	// write body
	buf.Write(bodyBuf.Bytes())

	resp := NewRpcResponse(&common.Header{}, buf)
	return resp
}

func mappingCode(code uint32) (esbCode string, message string) {
	switch code {
	case api.RouterUnavailableCode:
		esbCode, message = "B100", "no provider available(sidecar:404)."
	case api.NoHealthUpstreamCode:
		esbCode, message = "B100", "no health provider available(sidecar:502)."
	case api.TimeoutExceptionCode:
		esbCode, message = "B100", "invoke timeout(sidecar:504)."
	case api.CodecExceptionCode:
		esbCode, message = "B100", "decode error(sidecar:0)."
	default:
		esbCode, message = "B100", fmt.Sprintf("unknown error(sidecar:%d).", code)
	}

	return
}
