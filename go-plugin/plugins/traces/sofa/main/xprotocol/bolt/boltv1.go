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

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/plugins/trace/sofa/main/generator"
)

func init() {
	generator.RegisterDelegate(ProtocolName, Boltv1Delegate)
}

var (
	HeaderSofaRpcServiceKey                  = "service"
	ProtocolName            api.ProtocolName = "bolt" // protocol
)

func Boltv1Delegate(ctx context.Context, frame api.XFrame, span api.Span) {
	header := frame.GetHeader()
	lType, _ := config.GetListenerType(ctx)
	traceId, _ := header.Get(generator.TRACER_ID_KEY)
	if len(traceId) == 0 {
		span.SetTag(generator.SPAN_ID, "0")
		span.SetTag(generator.TRACE_ID, generator.IdGen().GenerateTraceId())
	} else {
		span.SetTag(generator.TRACE_ID, traceId)
		spanId, _ := header.Get(generator.RPC_ID_KEY)
		if lType == "ingress" {
			generator.AddSpanIdGenerator(generator.NewSpanIdGenerator(traceId, spanId))
		} else {
			span.SetTag(generator.PARENT_SPAN_ID, spanId)
			spanKey := &generator.SpanKey{TraceId: traceId, SpanId: spanId}
			if spanIdGenerator := generator.GetSpanIdGenerator(spanKey); spanIdGenerator != nil {
				spanId = spanIdGenerator.GenerateNextChildIndex()
			}
		}
		span.SetTag(generator.SPAN_ID, spanId)
		appName, _ := header.Get(generator.APP_NAME_KEY)
		span.SetTag(generator.CALLER_APP_NAME, appName)
	}
	method, _ := header.Get(generator.TARGET_METHOD_KEY)
	span.SetTag(generator.METHOD_NAME, method)
	span.SetTag(generator.PROTOCOL, "bolt")
	service, _ := header.Get(generator.SERVICE_KEY)
	span.SetTag(generator.SERVICE_NAME, service)
	bdata, _ := header.Get(generator.SOFA_TRACE_BAGGAGE_DATA)
	span.SetTag(generator.BAGGAGE_DATA, bdata)
	serviceKey, _ := header.Get(HeaderSofaRpcServiceKey)
	span.SetTag(generator.APP_SERVICE_NAME, fmt.Sprintf("%s@DEFAULT", serviceKey))
}
