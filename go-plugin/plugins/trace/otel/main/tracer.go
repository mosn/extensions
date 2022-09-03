package main

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel/propagation"
	otrace "go.opentelemetry.io/otel/trace"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/protocol/http"
)

const (
	HeadBaggageTraceID       = "Traceparent"
	HeadBaggageLetterTraceID = "traceparent"
)

var (
	tracerProvider otrace.TracerProvider
	_              api.Tracer = (*otelTracer)(nil)
)

type otelTracer struct {
	cfg            map[string]interface{}
	tracerProvider otrace.TracerProvider
}

// 单线程安全
func TracerBuilder(cfg map[string]interface{}) (api.Tracer, error) {
	tracerProvider, err := GetTracer(cfg)
	if err != nil {
		return nil, err
	}
	return &otelTracer{
		cfg:            cfg,
		tracerProvider: tracerProvider,
	}, nil
}

func (tracer *otelTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	switch req := request.(type) {
	case http.RequestHeader:
		return tracer.httpStart(ctx, req, startTime)
	case api.XFrame:
		return tracer.frameStart(ctx, req, startTime)
	}
	return trace.NewNooSpan()
}

func (tracer *otelTracer) httpStart(ctx context.Context, header http.RequestHeader, startTime time.Time) api.Span {
	omsg, ok := tracer.otelSpan(ctx, header)
	if !ok {
		return nil
	}
	return tracer.otel(ctx, header, omsg, startTime)
}

func (tracer *otelTracer) frameStart(ctx context.Context, xframe api.XFrame, startTime time.Time) api.Span {
	omsg, ok := tracer.otelSpan(ctx, xframe.GetHeader())
	if !ok {
		return nil
	}
	// ignore heartbeat
	if xframe.IsHeartbeatFrame() {
		return nil
	}
	return tracer.otel(ctx, xframe.GetHeader(), omsg, startTime)
}

func (tracer *otelTracer) otelSpan(ctx context.Context, header api.HeaderMap) (string, bool) {
	v, ok := header.Get(HeadBaggageTraceID)
	if ok {
		return v, ok
	}
	v, ok = header.Get(HeadBaggageLetterTraceID)
	if ok {
		return v, ok
	}
	enabled, _ := tracer.cfg["mosn_generator_span_enabled"].(string)
	if strings.EqualFold(enabled, "true") {
		return v, true
	}
	return "", false
}

func (tracer *otelTracer) otel(ctx context.Context, header api.HeaderMap, omsg string, startTime time.Time) api.Span {
	head := header.Clone()
	head.Set(HeadBaggageLetterTraceID, omsg)
	p := propagation.TraceContext{}
	pctx := p.Extract(ctx, &TextMapCarrier{head})
	ospan := NewSpan(ctx, startTime, tracer.tracerProvider)
	ospan.CreateLocalSpan(pctx)
	return ospan
}
