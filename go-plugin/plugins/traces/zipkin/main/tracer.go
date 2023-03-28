package main

import (
	"context"
	"strings"
	"time"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

var (
	tracerProvider  *zipkin.Tracer
	_               api.Tracer = (*zipkinTracer)(nil)
	DefaultSpanName            = "mosn"
)

type zipkinTracer struct {
	cfg            map[string]interface{}
	tracerProvider *zipkin.Tracer
	serviceName    string
}

func TracerBuilder(cfg map[string]interface{}) (api.Tracer, error) {
	tracerProvider, err := GetTracer(cfg)
	if err != nil {
		return nil, err
	}
	serviceName := cfg["service_name"].(string)
	return &zipkinTracer{
		cfg:            cfg,
		serviceName:    serviceName,
		tracerProvider: tracerProvider,
	}, nil
}

func (tracer *zipkinTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	switch req := request.(type) {
	case http.RequestHeader:
		return tracer.httpStart(ctx, req, startTime)
	case api.XFrame:
		return tracer.frameStart(ctx, req, startTime)
	}
	return trace.NewNooSpan()
}

func (tracer *zipkinTracer) httpStart(ctx context.Context, header http.RequestHeader, startTime time.Time) api.Span {
	zctx, err := tracer.zipkinSpan(ctx, header)
	if err != nil || zctx == nil {
		return nil
	}

	span := tracer.tracerProvider.StartSpan(
		getOperationName(header.RequestURI()),
		zipkin.Parent(*zctx),
		zipkin.Kind(tracer.getType(ctx)),
		zipkin.StartTime(startTime),
	)
	// getLocalHostPort get host and port from context
	return tracer.zipkin(ctx, header, span, startTime)
}

func (tracer *zipkinTracer) frameStart(ctx context.Context, xframe api.XFrame, startTime time.Time) api.Span {
	zctx, err := tracer.zipkinSpan(ctx, xframe.GetHeader())
	if err != nil || zctx == nil {
		return nil
	}
	// ignore heartbeat
	if xframe.IsHeartbeatFrame() {
		return nil
	}
	span := tracer.tracerProvider.StartSpan(DefaultSpanName,
		zipkin.Parent(*zctx),
		zipkin.Kind(tracer.getType(ctx)),
		zipkin.StartTime(startTime),
	)
	return tracer.zipkin(ctx, xframe.GetHeader(), span, startTime)
}

func (tracer *zipkinTracer) zipkinSpan(bctx context.Context, header api.HeaderMap) (*model.SpanContext, error) {
	if ctx, _ := tracer.parseSpanFromHeader(bctx, header); ctx != nil {
		return ctx, nil
	}
	if ctx, _ := tracer.parseSpanFromHeader2(bctx, header); ctx != nil {
		return ctx, nil
	}
	enabled, _ := tracer.cfg["mosn_generator_span_enabled"].(string)
	if strings.EqualFold(enabled, "true") {
		ctx, err := b3.ParseHeaders("", "", "", "1", "0")
		if err != nil {
			return nil, err
		}
		return ctx, nil
	}
	return nil, nil
}

func (tracer *zipkinTracer) parseSpanFromHeader(bctx context.Context, header api.HeaderMap) (*model.SpanContext, error) {
	singleHeader, ok := header.Get(b3.Context)
	if !ok {
		return nil, nil
	}
	ctx, err := b3.ParseSingleHeader(singleHeader)
	if err != nil {
		log.DefaultLogger.Errorf("[zipkin][parseSpanFromHeader]get span failed,err:%s header:%s", err, singleHeader)
		return nil, err
	}
	return ctx, nil
}

func (tracer *zipkinTracer) parseSpanFromHeader2(bctx context.Context, header api.HeaderMap) (*model.SpanContext, error) {
	var (
		hdrTraceID, ok     = header.Get(b3.TraceID)
		hdrSpanID, _       = header.Get(b3.SpanID)
		hdrParentSpanID, _ = header.Get(b3.ParentSpanID)
		hdrSampled, _      = header.Get(b3.Sampled)
		hdrFlgs, _         = header.Get(b3.Flags)
	)
	if !ok {
		return nil, nil
	}
	ctx, err := b3.ParseHeaders(hdrTraceID, hdrSpanID, hdrParentSpanID, hdrSampled, hdrFlgs)
	if err != nil {
		log.DefaultLogger.Errorf("[zipkin][parseSpanFromHeader2]get span failed,err:%s", err)
		return nil, err
	}
	return ctx, nil
}

func (tracer *zipkinTracer) getType(ctx context.Context) model.Kind {
	if ltype, ok := config.GetListenerType(ctx); ok && ltype == "ingress" {
		return model.Server
	}
	return model.Client
}

func (tracer *zipkinTracer) zipkin(ctx context.Context, header api.HeaderMap, span zipkin.Span, startTime time.Time) api.Span {
	ospan := NewSpan(ctx, startTime, tracer.tracerProvider)
	ospan.CreateLocalSpan(span)
	return ospan
}
