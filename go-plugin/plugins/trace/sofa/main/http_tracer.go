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

package main

import (
	"context"
	"time"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/extensions/go-plugin/plugins/trace/sofa/main/generator"
	"mosn.io/pkg/protocol/http"
)

type HTTPTracer struct {
	tracer api.Tracer
	cfg    map[string]interface{}
}

func NewHTTPTracer(config map[string]interface{}) (api.Tracer, error) {
	config["server_name"] = "springcloud-server-digest.log"
	config["client_name"] = "springcloud-client-digest.log"
	tracer, err := NewTracer(config)
	if err != nil {
		return nil, err
	}
	return &HTTPTracer{
		tracer: tracer,
		cfg:    config,
	}, nil
}

func (t *HTTPTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	span := t.tracer.Start(ctx, request, startTime)
	header, ok := request.(http.RequestHeader)
	if !ok || header.RequestHeader == nil {
		return span
	}
	t.HTTPDelegate(ctx, header, span)
	if len(span.Tag(generator.TRACE_ID)) != 0 {
		return span
	}
	return trace.NewNooSpan()
}

func (t *HTTPTracer) HTTPDelegate(ctx context.Context, header http.RequestHeader, span api.Span) {
	lType, _ := config.GetListenerType(ctx)
	traceId, _ := header.Get(generator.HTTP_TRACER_ID_KEY)
	if len(traceId) == 0 {
		traceId = generator.IdGen().GenerateTraceId()
		span.SetTag(generator.TRACE_ID, traceId)
		span.SetTag(generator.SPAN_ID, "0")
	} else {
		span.SetTag(generator.TRACE_ID, traceId)
		spanId, _ := header.Get(generator.HTTP_RPC_ID_KEY)
		if lType == "INGRESS" {
			generator.AddSpanIdGenerator(generator.NewSpanIdGenerator(traceId, spanId))
		} else if lType == "EGRESS" {
			span.SetTag(generator.PARENT_SPAN_ID, spanId)
			spanKey := &generator.SpanKey{TraceId: traceId, SpanId: spanId}
			if spanIdGenerator := generator.GetSpanIdGenerator(spanKey); spanIdGenerator != nil {
				spanId = spanIdGenerator.GenerateNextChildIndex()
			}
		}
		span.SetTag(generator.SPAN_ID, spanId)
	}
	span.SetTag(generator.REQUEST_URL, string(header.RequestURI()))
	if lType == "EGRESS" {
		span.SetTag(generator.CALLER_APP_NAME, string(header.Peek(generator.APP_NAME_KEY)))
	}
	span.SetTag(generator.METHOD_NAME, string(header.Peek(generator.TARGET_METHOD_KEY)))
	span.SetTag(generator.SERVICE_NAME, string(header.Peek(generator.SERVICE_KEY)))
	span.SetTag(generator.BAGGAGE_DATA, string(header.Peek(generator.SOFA_TRACE_BAGGAGE_DATA)))
	span.SetTag(generator.PROTOCOL_FRAME, "HTTP")
}
