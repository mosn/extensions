package main

import (
	"context"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	language_agent "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

const (
	HeadBaggageTraceID       = "Sw8"
	HeadBaggageLetterTraceID = "sw8"
	MOSNComponentID          = 5003
)

var (
	_ api.Tracer = (*skyTracer)(nil)
)

type skyTracer struct {
	cfg            map[string]interface{}
	tracerProvider *go2sky.Tracer
}

func TracerBuilder(cfg map[string]interface{}) (api.Tracer, error) {
	tracerProvider, err := GetTracer(cfg)
	if err != nil {
		return nil, err
	}
	return &skyTracer{
		cfg:            cfg,
		tracerProvider: tracerProvider,
	}, nil
}

func (tracer *skyTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	switch req := request.(type) {
	case http.RequestHeader:
		return tracer.httpStart(ctx, req, startTime)
	case api.XFrame:
		return tracer.frameStart(ctx, req, startTime)
	}
	return trace.NewNooSpan()
}

func (tracer *skyTracer) httpStart(ctx context.Context, header http.RequestHeader, startTime time.Time) api.Span {
	omsg, ok := tracer.skySpan(ctx, header)
	if !ok {
		return nil
	}
	requestURI := string(header.RequestURI())
	entry, nctx, err := tracer.tracerProvider.CreateEntrySpan(ctx, requestURI, func() (string, error) {
		return omsg, nil
	})
	if err != nil {
		log.DefaultLogger.Errorf("[SkyWalking] [tracer] [http1] create entry span error, err: %v", err)
		return nil
	}
	entry.SetSpanLayer(language_agent.SpanLayer_Http)
	entry.SetComponent(MOSNComponentID)
	// new span
	ospan := NewSpan(ctx, startTime, tracer.tracerProvider)
	ospan.CreateLocalHttpSpan(nctx, header, entry)
	return ospan
}

func (tracer *skyTracer) frameStart(ctx context.Context, xframe api.XFrame, startTime time.Time) api.Span {
	omsg, ok := tracer.skySpan(ctx, xframe.GetHeader())
	if !ok {
		return nil
	}
	// ignore heartbeat
	if xframe.IsHeartbeatFrame() {
		return nil
	}
	entry, nctx, err := tracer.tracerProvider.CreateEntrySpan(ctx, "mosn", func() (string, error) {
		return omsg, nil
	})
	if err != nil {
		log.DefaultLogger.Errorf("[SkyWalking] [tracer] [http1] create entry span error, err: %v", err)
		return nil
	}
	entry.SetSpanLayer(language_agent.SpanLayer_RPCFramework)
	entry.SetComponent(MOSNComponentID)
	// new span
	ospan := NewSpan(ctx, startTime, tracer.tracerProvider)
	ospan.CreateLocalRpcSpan(nctx, entry)
	return ospan
}

func (tracer *skyTracer) skySpan(ctx context.Context, header api.HeaderMap) (string, bool) {
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
		// TODO traceid
		return v, true
	}
	return "", false
}
