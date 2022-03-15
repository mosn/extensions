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
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/extensions/go-plugin/pkg/common/safe"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
	"time"
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
		return nil, fmt.Errorf("decode cd rpc Error, unknown request type = %s", frameType)
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

	// decode app-header
	v, err := parseXmlHeader(req.Payload.Bytes(), startAppHeader, endAppHeader)
	if err != nil {
		// should never happen
		log.DefaultLogger.Errorf("failed to resolve cd proto app header, err %v, data: %s", err, string(req.Payload.Bytes()))
	}

	// decode app header
	var (
		branchId string
		userId   string
	)

	if len(v.WrapData) > 0 {
		for _, d := range v.WrapData {
			if d.Field != nil { // plain field
				switch d.Name {
				case branchIdKey:
					branchId = d.Field.Value
				case userIdKey:
					userId = d.Field.Value
				}
			}
		}
	} else {
		// should never happen
		log.DefaultLogger.Warnf("resolved empty cd proto app header, data: %s", string(req.Payload.Bytes()))
	}

	if branchId == "" {
		branchId, _ = req.Get(branchIdKey)
	}

	if userId == "" {
		userId, _ = req.Get(userIdKey)
	}

	// 1. write xml header
	bodyBuf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")

	// 2. write service
	bodyBuf.WriteString("<service>")
	{
		bodyBuf.WriteString("<sys-header>")
		{
			bodyBuf.WriteString("<data name=\"SERVICE_CODE\">")
			code, _ := req.Get(serviceCodeKey)
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(code), code))
			bodyBuf.WriteString("</data>")
		}
		{
			bodyBuf.WriteString("<data name=\"SERVICE_SCENE\">")
			scene, _ := req.Get(serviceSceneKey)
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(scene), scene))
			bodyBuf.WriteString("</data>")
		}
		{
			bodyBuf.WriteString("<data name=\"CONSUMER_ID\">")
			consumerId, _ := req.Get(consumerIdKey)
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(consumerId), consumerId))
			bodyBuf.WriteString("</data>")
		}
		{
			bodyBuf.WriteString("<data name=\"CONSUMER_SEQ_NO\">")
			consumerSeqNo, _ := req.Get(consumerSeqNoKey)
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(consumerSeqNo), consumerSeqNo))
			bodyBuf.WriteString("</data>")
		}
		{
			dateLayout := "20060102"
			timeLayout := "150405"
			now := time.Now()

			dateVal := now.Format(dateLayout)
			timeVal := now.Format(timeLayout)
			bodyBuf.WriteString("<data name=\"TRAN_DATE\">")
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(dateVal), dateVal))
			bodyBuf.WriteString("</data>")

			bodyBuf.WriteString("<data name=\"TRAN_TIMESTAMP\">")
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(timeVal), timeVal))
			bodyBuf.WriteString("</data>")
		}
		{
			bodyBuf.WriteString("<data name=\"RET_STATUS\">")
			bodyBuf.WriteString("<field length=\"1\" scale=\"0\" type=\"string\">F</field>") // failed
			bodyBuf.WriteString("</data>")

			bodyBuf.WriteString("<data name=\"RET\">")
			{
				bodyBuf.WriteString("<array><struct>")
				{
					bodyBuf.WriteString("<data name=\"RET_CODE\">")

					code, message := mappingCode(statusCode)
					bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(code), code)) // failed
					bodyBuf.WriteString("</data>")

					bodyBuf.WriteString("<data name=\"RET_MSG\">")
					bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(message), message)) // failed
					bodyBuf.WriteString("</data>")
				}
				bodyBuf.WriteString("</struct></array>")
			}
			bodyBuf.WriteString("</data>")
		}
		{
			bodyBuf.WriteString("<data name=\"TRAN_ID\">")
			tranId, _ := req.Get(tranIdKey)
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(tranId), tranId))
			bodyBuf.WriteString("</data>")
		}
		bodyBuf.WriteString("</sys-header>")
	}

	// 3. write appHeader
	{
		bodyBuf.WriteString("<app-header>")
		{
			bodyBuf.WriteString("<data name=\"BRANCH_ID\">")
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(branchId), branchId))
			bodyBuf.WriteString("</data>")

			bodyBuf.WriteString("<data name=\"USER_ID\">")
			bodyBuf.WriteString(fmt.Sprintf("<field length=\"%d\" scale=\"0\" type=\"string\">%s</field>", len(userId), userId))
			bodyBuf.WriteString("</data>")
		}
		bodyBuf.WriteString("</app-header>")
	}

	// 4. write local-header
	{
		bodyBuf.WriteString("<local-header>")
		{
			bodyBuf.WriteString("<data name=\"LOCAL_HEAD\">")
			{
				bodyBuf.WriteString("<struct/>")
			}
			bodyBuf.WriteString("</data>")
		}
		bodyBuf.WriteString("</local-header>")
	}

	// 5. write body
	bodyBuf.WriteString("<body/>")
	bodyBuf.WriteString("</service>")

	// 10 byte length + string body
	buf := buffer.GetIoBuffer(10 + bodyBuf.Len())

	// write 10 byte length + xml body
	proto.prefixOfZero(buf, bodyBuf.Len())
	buf.Write(bodyBuf.Bytes())

	resp := NewRpcResponse(&common.Header{}, buf)
	return resp
}

func mappingCode(code uint32) (esbCode string, message string) {
	switch code {
	case api.RouterUnavailableCode:
		esbCode, message = "999999", "no provider available(sidecar:404)."
	case api.NoHealthUpstreamCode:
		esbCode, message = "999999", "no health provider available(sidecar:502)."
	case api.TimeoutExceptionCode:
		esbCode, message = "999999", "invoke timeout(sidecar:504)."
	case api.CodecExceptionCode:
		esbCode, message = "999999", "decode error(sidecar:0)."
	default:
		esbCode, message = "999999", fmt.Sprintf("unknown error(sidecar:%d).", code)
	}

	return
}
