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
	"mosn.io/extensions/go-plugin/plugins/traces/sofa/main/generator"
	"mosn.io/pkg/log"
)

type XTracer struct {
	cfg    map[string]interface{}
	tracer api.Tracer
}

func NewRpcTracer(config map[string]interface{}) (api.Tracer, error) {
	config["server_name"] = "rpc-server-digest.log"
	config["client_name"] = "rpc-client-digest.log"
	config["tracer_type"] = "rpc"
	tracer, err := NewTracer(config)
	if err != nil {
		return nil, err
	}
	return &XTracer{
		tracer: tracer,
		cfg:    config,
	}, nil
}

func (t *XTracer) Start(ctx context.Context, frame interface{}, startTime time.Time) api.Span {
	span := t.tracer.Start(ctx, frame, startTime)
	xframe, ok := frame.(api.XFrame)
	if !ok || xframe == nil {
		return span
	}
	// ignore heartbeat
	if xframe.IsHeartbeatFrame() {
		return span
	}
	proto, _ := config.GetDownstreamProtocol(ctx)
	if delegate := generator.GetDelegate(api.ProtocolName(proto)); delegate != nil {
		delegate(ctx, xframe, span)
	}
	if len(span.Tag(generator.TRACE_ID)) != 0 {
		return span
	}
	span.SetTag(generator.TRACE_ID, generator.IdGen().GenerateTraceId())
	span.SetTag(generator.SPAN_ID, "0")
	log.DefaultLogger.Warnf("[cloud_sofa] [tracer] the %s of span not found ", proto)
	return span
}
