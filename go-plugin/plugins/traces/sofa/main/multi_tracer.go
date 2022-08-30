package main

import (
	"context"
	"time"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/protocol/http"
)

type MultiTracer struct {
	htracer api.Tracer
	rtracer api.Tracer
}

func TracerBuilder(cfg map[string]interface{}) (api.Tracer, error) {
	htracer, err := NewHTTPTracer(cfg)
	if err != nil {
		return nil, err
	}
	rtracer, err := NewRpcTracer(cfg)
	if err != nil {
		return nil, err
	}
	return &MultiTracer{
		htracer: htracer,
		rtracer: rtracer,
	}, nil
}

func (t *MultiTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	switch req := request.(type) {
	case http.RequestHeader:
		return t.htracer.Start(ctx, req, startTime)
	case api.XFrame:
		return t.rtracer.Start(ctx, req, startTime)
	}
	return trace.NewNooSpan()
}
