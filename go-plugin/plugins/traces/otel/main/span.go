package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otrace "go.opentelemetry.io/otel/trace"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

const (
	instrumentationName = "gitlab.alipay-inc.com/ant-mesh/mosn"
)

type OtelSpan struct {
	trace.NoopSpan
	tid       string
	sid       string
	psid      string
	startTime time.Time
	ctx       context.Context

	operationName string
	kvs           []attribute.KeyValue
	provider      otrace.TracerProvider
	entrySpan     otrace.Span
	entryctx      context.Context
}

func NewSpan(ctx context.Context, startTime time.Time, provider otrace.TracerProvider) *OtelSpan {
	h := &OtelSpan{
		startTime: startTime,
		ctx:       ctx,
		provider:  provider,
		kvs:       make([]attribute.KeyValue, 0, 20),
	}
	return h
}

func (h *OtelSpan) TraceId() string {
	return h.tid
}

func (h *OtelSpan) SpanId() string {
	return h.sid
}

func (h *OtelSpan) ParentSpanId() string {
	return h.psid
}

func (h *OtelSpan) InjectContext(headers api.HeaderMap, reqInfo api.RequestInfo) {
	host := reqInfo.UpstreamLocalAddress()
	if header, ok := headers.(http.RequestHeader); ok {
		uri := string(header.RequestURI())
		url := strings.Join([]string{"http://", host, uri}, "")
		h.kvs = append(h.kvs, semconv.HTTPURLKey.String(url))
		h.kvs = append(h.kvs, semconv.HTTPMethodKey.String(string(header.Method())))
		h.operationName = uri
	}
	p := propagation.TraceContext{}
	p.Inject(h.entryctx, &TextMapCarrier{headers})
}

func (h *OtelSpan) SetRequestInfo(reqInfo api.RequestInfo) {
	if host := reqInfo.UpstreamHost(); host != nil {
		h.kvs = append(h.kvs, attribute.Key("upstream.address").String(host.AddressString()))
		semconv.HostNameKey.String(host.Hostname())
	}
	if addr := reqInfo.DownstreamRemoteAddress(); addr != nil {
		h.kvs = append(h.kvs, attribute.Key("downstream.address").String(addr.String()))
	}
	h.kvs = append(h.kvs, attribute.Key("request.size").Int64(int64(reqInfo.BytesReceived())))
	h.kvs = append(h.kvs, attribute.Key("respone.size").Int64(int64(reqInfo.BytesSent())))
	h.kvs = append(h.kvs, attribute.Key("duration").Int64(reqInfo.Duration().Nanoseconds()))
	process := reqInfo.ProcessTimeDuration().Nanoseconds()
	if process == 0 {
		process = reqInfo.RequestFinishedDuration().Nanoseconds()
	}
	h.kvs = append(h.kvs, attribute.Key("mosn.process.duration").Int64(process))
	h.kvs = append(h.kvs, attribute.Key("mosn.process.request.duration").Int64(reqInfo.RequestFinishedDuration().Nanoseconds()))
	h.kvs = append(h.kvs, attribute.Key("mosn.process.respone.duration").Int64(reqInfo.ResponseReceivedDuration().Nanoseconds()))
	if reqInfo.ResponseCode() != api.SuccessCode {
		ok := reqInfo.GetResponseFlag(trace.MosnProcessFailedFlags)
		h.kvs = append(h.kvs, attribute.Key("mosn.process.fail").Bool(ok))
		h.entrySpan.SetStatus(codes.Error, "reqeust failed")
		if ok {
			h.entrySpan.RecordError(errors.New("mosn error"))
		} else {
			h.entrySpan.RecordError(errors.New("biz error"))
		}
	} else {
		h.entrySpan.SetStatus(codes.Ok, "")
	}
}

func (h *OtelSpan) FinishSpan() {
	if h.entrySpan != nil {
		h.ParseVariable(h.ctx)
		h.entrySpan.SetName(h.operationName)
		h.entrySpan.SetAttributes(h.kvs...)
		h.entrySpan.End()
		h.log()
	}
}

func (h *OtelSpan) SetOperation(operation string) {
	h.operationName = operation
}

func (h *OtelSpan) log() {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		kvs, _ := json.Marshal(h.kvs)
		log.DefaultLogger.Debugf("trace:%s pid:%s parentid:%s operationName:%s,kvs:%s", h.tid, h.sid, h.psid, h.operationName, kvs)
	}
}

// otrace.WithSpanKind(h.ctx.SpanKind())
func (h *OtelSpan) CreateLocalSpan(ctx context.Context) {
	span := otrace.SpanFromContext(ctx)
	var opts []otrace.SpanStartOption
	opts = append(opts, otrace.WithTimestamp(h.startTime))
	opts = append(opts, otrace.WithSpanKind(h.SpanKind()))
	h.entryctx, h.entrySpan = h.provider.Tracer(instrumentationName).Start(ctx, h.operationName, opts...)
	h.psid = span.SpanContext().SpanID().String()
	h.sid = h.entrySpan.SpanContext().SpanID().String()
	h.tid = h.entrySpan.SpanContext().TraceID().String()
}

func (h *OtelSpan) SpanKind() otrace.SpanKind {
	if ltype, ok := config.GetListenerType(h.ctx); ok && ltype == "ingress" {
		return otrace.SpanKindServer
	}
	return otrace.SpanKindClient
}

func (h *OtelSpan) ParseVariable(ctx context.Context) {
	dp, _ := config.GetDownstreamProtocol(h.ctx)
	if len(dp) != 0 {
		h.kvs = append(h.kvs, attribute.Key("downstream.protocol").String(string(dp)))
	}
	up, _ := config.GetUpstreamProtocol(h.ctx)
	if len(up) != 0 {
		h.kvs = append(h.kvs, attribute.Key("upstream.protocol").String(string(up)))
	} else {
		h.kvs = append(h.kvs, attribute.Key("upstream.protocol").String(string(dp)))
	}
}
