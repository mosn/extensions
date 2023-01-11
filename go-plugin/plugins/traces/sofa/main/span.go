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
	"fmt"
	"strconv"
	"time"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/extensions/go-plugin/plugins/traces/sofa/main/generator"
	"mosn.io/pkg/log"
)

type SofaRPCSpan struct {
	ctx           context.Context
	ingressLogger *log.Logger
	egressLogger  *log.Logger
	startTime     time.Time
	endTime       time.Time
	tags          [generator.TRACE_END]string
	traceId       string
	spanId        string
	parentSpanId  string
	operationName string
	appName       string
	pod           bool
	cluster       string
}

func (s *SofaRPCSpan) TraceId() string {
	return s.traceId
}

func (s *SofaRPCSpan) SpanId() string {
	return s.spanId
}

func (s *SofaRPCSpan) ParentSpanId() string {
	return s.parentSpanId
}

func (s *SofaRPCSpan) SetOperation(operation string) {
	s.operationName = operation
}

func (s *SofaRPCSpan) SetTag(key uint64, value string) {
	if key == generator.TRACE_ID {
		s.traceId = value
	} else if key == generator.SPAN_ID {
		s.spanId = value
	} else if key == generator.PARENT_SPAN_ID {
		s.parentSpanId = value
	}

	s.tags[key] = value
}

func (s *SofaRPCSpan) SetRequestInfo(reqinfo api.RequestInfo) {
	s.tags[generator.REQUEST_SIZE] = strconv.FormatInt(int64(reqinfo.BytesReceived()), 10) + "B"
	s.tags[generator.RESPONSE_SIZE] = strconv.FormatInt(int64(reqinfo.BytesSent()), 10) + "B"
	if reqinfo.UpstreamHost() != nil {
		s.tags[generator.UPSTREAM_HOST_ADDRESS] = reqinfo.UpstreamHost().AddressString()
	}
	if reqinfo.DownstreamRemoteAddress() != nil {
		s.tags[generator.DOWNSTEAM_HOST_ADDRESS] = reqinfo.DownstreamRemoteAddress().String()
	}
	s.tags[generator.RESULT_STATUS] = strconv.Itoa(reqinfo.ResponseCode())
	s.tags[generator.MOSN_PROCESS_TIME] = reqinfo.ProcessTimeDuration().String()
	s.tags[generator.MOSN_PROCESS_FAIL] = strconv.FormatBool(reqinfo.GetResponseFlag(trace.MosnProcessFailedFlags))
	s.tags[generator.DURATION] = reqinfo.Duration().String()
	s.tags[generator.REQUEST_DURATION] = reqinfo.RequestFinishedDuration().String()
	s.tags[generator.RESPONSE_DURATION] = reqinfo.ResponseReceivedDuration().String()
	s.tags[generator.UPSTREAM_DURATION] = strconv.Itoa(int(reqinfo.Duration().Milliseconds()-reqinfo.ProcessTimeDuration().Milliseconds())) + "ms"
}

func (s *SofaRPCSpan) Tag(key uint64) string {
	return s.tags[key]
}

func (s *SofaRPCSpan) FinishSpan() {
	s.endTime = time.Now()
	err := s.log()
	if err != nil {
		log.DefaultLogger.Warnf("Channel is full, discard span, trace id is " + s.traceId + ", span id is " + s.spanId)
	}
}

func (s *SofaRPCSpan) InjectContext(requestHeaders api.HeaderMap, requestInfo api.RequestInfo) {
}

func (s *SofaRPCSpan) SpawnChild(operationName string, startTime time.Time) api.Span {
	return nil
}

func (s *SofaRPCSpan) SetStartTime(startTime time.Time) {
	s.startTime = startTime
}

func (s *SofaRPCSpan) String() string {
	return fmt.Sprintf("TraceId:%s;SpanId:%s;Duration:%s;ProtocolName:%s;ServiceName:%s;requestSize:%s;responseSize:%s;upstreamHostAddress:%s;downstreamRemoteHostAdress:%s",
		s.tags[generator.TRACE_ID],
		s.tags[generator.SPAN_ID],
		strconv.FormatInt(s.endTime.Sub(s.startTime).Nanoseconds()/1000000, 10),
		s.tags[generator.PROTOCOL],
		s.tags[generator.SERVICE_NAME],
		s.tags[generator.REQUEST_SIZE],
		s.tags[generator.RESPONSE_SIZE],
		s.tags[generator.UPSTREAM_HOST_ADDRESS],
		s.tags[generator.DOWNSTEAM_HOST_ADDRESS])
}

func (s *SofaRPCSpan) EndTime() time.Time {
	return s.endTime
}

func (s *SofaRPCSpan) StartTime() time.Time {
	return s.startTime
}

func (s *SofaRPCSpan) parseVariable(ctx context.Context) {
	/*
		if methodName, _ := variable.GetString(ctx, "x-mosn-method"); len(methodName) != 0 {
			s.SetTag(generator.METHOD_NAME, methodName)
		}
		if appName, _ := variable.GetString(ctx, "x-mosn-caller-app"); len(appName) != 0 {
			s.SetTag(generator.CALLER_APP_NAME, appName)
		}
		if service, _ := variable.GetString(ctx, "x-mosn-target-app"); len(service) != 0 {
			s.SetTag(generator.TARGET_APP_NAME, service)
		}
		if dataId, _ := variable.GetString(ctx, "x-mosn-data-id"); len(dataId) != 0 {
			s.SetTag(generator.APP_SERVICE_NAME, dataId)
			if len(s.operationName) == 0 {
				s.operationName = dataId
			}
		}

	*/
	dp, _ := config.GetDownstreamProtocol(ctx)
	if len(dp) != 0 {
		s.SetTag(generator.DOWNSTREAM_PROTOCOL, string(dp))
		s.SetTag(generator.PROTOCOL, string(dp))
	}
	up, _ := config.GetUpstreamProtocol(ctx)
	if len(up) != 0 {
		s.SetTag(generator.UPSTREAM_PROTOCOL, string(up))
	} else {
		s.SetTag(generator.UPSTREAM_PROTOCOL, string(dp))
	}
	if ltype, ok := config.GetListenerType(ctx); ok && ltype == "ingress" {
		s.SetTag(generator.SPAN_TYPE, "ingress")
	} else {
		s.SetTag(generator.SPAN_TYPE, "egress")
	}
}

func (s *SofaRPCSpan) log() error {
	s.parseVariable(s.ctx)
	s.baggage()
	if s.tags[generator.PROTOCOL_FRAME] != "rpc" {
		if s.tags[generator.SPAN_TYPE] == "ingress" {
			return s.serverHttpLogger()
		}
		if s.tags[generator.SPAN_TYPE] == "egress" {
			return s.clientHttpLogger()
		}
	} else {
		if s.tags[generator.SPAN_TYPE] == "ingress" {
			return s.serverRpcLogger()
		}
		if s.tags[generator.SPAN_TYPE] == "egress" {
			return s.clientRpcLogger()
		}
	}
	return nil
}

func (s *SofaRPCSpan) baggage() {
	penMap := s.Tag(generator.BAGGAGE_DATA)
	appService := s.Tag(generator.SERVICE_NAME)
	appServiceVersion := s.Tag(generator.APP_SERVICE_VERSION)
	penMap = s.joinAppServiceForCloud(penMap, appService, appServiceVersion)
	s.SetTag(generator.BAGGAGE_DATA, penMap)
}

func (s *SofaRPCSpan) joinAppServiceForCloud(origin, appService, appServiceVersion string) string {
	return origin + "mosn_cluster=" + s.cluster + "&mosn_data_id=" + appService + "&mosn_data_ver=" + appServiceVersion + "&"
}

func NewSpan(ctx context.Context, startTime time.Time) *SofaRPCSpan {
	return &SofaRPCSpan{
		ctx:       ctx,
		startTime: startTime,
	}
}
