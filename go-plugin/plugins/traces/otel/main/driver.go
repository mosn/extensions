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
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otrace "go.opentelemetry.io/otel/trace"
	"mosn.io/api"
	"mosn.io/pkg/log"
)

const (
	OtelDriverName = "OTEL"
)

var (
	Version = ""
)

// 单线程安全
func GetTracer(config map[string]interface{}) (otrace.TracerProvider, error) {
	var err error
	if tracerProvider == nil {
		tracerProvider, err = newTracer(config)
	}
	return tracerProvider, err
}

func otelAgentClient(config map[string]interface{}) (otlptrace.Client, error) {
	otelAgentAddr, ok := config["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address is empty")
	}
	if stype, ok := config["service_type"].(string); ok && strings.EqualFold(stype, "rpc") {
		// rpc client
		return otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(otelAgentAddr)), nil
	}
	opts := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(time.Second * 3),
	}
	if strings.Contains(otelAgentAddr, "//") {
		httpurl, err := url.Parse(otelAgentAddr)
		if err != nil {
			return nil, err
		}
		opts = append(opts, otlptracehttp.WithEndpoint(httpurl.Host))
		if len(httpurl.Path) != 0 {
			opts = append(opts, otlptracehttp.WithURLPath(httpurl.Path))
		}
	} else {
		opts = append(opts, otlptracehttp.WithEndpoint(otelAgentAddr))
	}
	return otlptracehttp.NewClient(opts...), nil
}

func newTracer(config map[string]interface{}) (otrace.TracerProvider, error) {
	otel.SetErrorHandler(TextMapCarrier{})
	traceClient, err := otelAgentClient(config)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	service, _ := config["service"].(string)
	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceVersionKey.String(Version),
			semconv.ServiceNameKey.String(service),
		),
	)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider, nil
}

type TextMapCarrier struct {
	api.HeaderMap
}

func (tm TextMapCarrier) Get(key string) string {
	v, _ := tm.HeaderMap.Get(key)
	return v
}

func (tm TextMapCarrier) Keys() []string {
	var keys []string
	tm.HeaderMap.Range(
		func(key, value string) bool {
			keys = append(keys, key)
			return true
		})
	return keys
}

func (tm TextMapCarrier) Handle(err error) {
	log.DefaultLogger.Errorf("[tracer][otel] err:%s", err)
}
