package main

import (
	"context"
	"encoding/json"
	"mosn.io/extensions/go-plugin/pkg/keys"
	"mosn.io/pkg/variable"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

type kv struct {
	Key   string
	Value string
}

type ZipkinSpan struct {
	trace.NoopSpan
	tid       string
	sid       string
	psid      string
	startTime time.Time
	ctx       context.Context

	operationName string
	kvs           []kv
	provider      *zipkin.Tracer
	entrySpan     zipkin.Span
}

func NewSpan(ctx context.Context, startTime time.Time, provider *zipkin.Tracer) *ZipkinSpan {
	h := &ZipkinSpan{
		startTime: startTime,
		ctx:       ctx,
		provider:  provider,
		kvs:       make([]kv, 0, 10),
	}
	return h
}

func (h *ZipkinSpan) TraceId() string {
	return h.tid
}

func (h *ZipkinSpan) SpanId() string {
	return h.sid
}

func (h *ZipkinSpan) ParentSpanId() string {
	return h.psid
}

func (h *ZipkinSpan) InjectContext(headers api.HeaderMap, reqInfo api.RequestInfo) {
	if header, ok := headers.(http.RequestHeader); ok {
		host := reqInfo.UpstreamLocalAddress()
		requestURI := string(header.RequestURI())
		url := strings.Join([]string{"http://", host, requestURI}, "")
		h.kvs = append(h.kvs, kv{"http.url", url})
		h.kvs = append(h.kvs, kv{"http.method", string(header.Method())})
		h.operationName = requestURI
	}
	trace := b3.BuildSingleHeader(h.entrySpan.Context())
	headers.Set(b3.Context, trace)
}

func (h *ZipkinSpan) SetRequestInfo(reqInfo api.RequestInfo) {
	if host := reqInfo.UpstreamHost(); host != nil {
		h.kvs = append(h.kvs, kv{"upstream.address", host.AddressString()})
		addr, _, err := net.SplitHostPort(host.AddressString())
		if err == nil {
			endpoint := &model.Endpoint{
				ServiceName: host.Hostname(),
				IPv4:        net.ParseIP(addr),
			}
			h.entrySpan.SetRemoteEndpoint(endpoint)
		}
	}
	if addr := reqInfo.DownstreamRemoteAddress(); addr != nil {
		h.kvs = append(h.kvs, kv{"downstream.address", addr.String()})
	}
	h.kvs = append(h.kvs, kv{"request.size", strconv.Itoa(int(reqInfo.BytesReceived()))})
	h.kvs = append(h.kvs, kv{"response.size", strconv.Itoa(int(reqInfo.BytesSent()))})
	h.kvs = append(h.kvs, kv{"duration", strconv.Itoa(int(reqInfo.Duration().Nanoseconds()))})
	process := reqInfo.ProcessTimeDuration().Nanoseconds()
	if process == 0 {
		process = reqInfo.RequestFinishedDuration().Nanoseconds()
	}
	h.kvs = append(h.kvs, kv{"mosn.process.duration", strconv.Itoa(int(process))})
	h.kvs = append(h.kvs, kv{"mosn.process.request.duration", strconv.Itoa(int(reqInfo.RequestFinishedDuration().Nanoseconds()))})
	h.kvs = append(h.kvs, kv{"mosn.process.response.duration", strconv.Itoa(int(reqInfo.ResponseReceivedDuration().Nanoseconds()))})
	if reqInfo.ResponseCode() != api.SuccessCode {
		ok := reqInfo.GetResponseFlag(trace.ProcessFailedFlags)
		h.kvs = append(h.kvs, kv{"mosn.process.fail", strconv.FormatBool(ok)})
	}
}

func (h *ZipkinSpan) FinishSpan() {
	if h.entrySpan != nil {
		kvs := h.ParseVariable(h.ctx)
		h.entrySpan.SetName(h.operationName)
		h.addTag(h.entrySpan, kvs)
		h.entrySpan.Finish()
		h.log(kvs)
	}
}

func (h *ZipkinSpan) addTag(span zipkin.Span, kvs []kv) {
	for _, kv := range kvs {
		span.Tag(kv.Key, kv.Value)
	}
}

func (h *ZipkinSpan) SetOperation(operation string) {
	h.operationName = operation
}

func (h *ZipkinSpan) log(kvs []kv) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		kvs, _ := json.Marshal(h.kvs)
		log.DefaultLogger.Debugf("trace:%s pid:%s parentid:%s operationName:%s,kvs:%s", h.tid, h.sid, h.psid, h.operationName, kvs)
	}
}

func (h *ZipkinSpan) CreateLocalSpan(span zipkin.Span) {
	h.entrySpan = span
	if pid := h.entrySpan.Context().ParentID; pid != nil {
		h.psid = pid.String()
	}
	h.tid = h.entrySpan.Context().TraceID.String()
	h.sid = h.entrySpan.Context().ID.String()
}

func (h *ZipkinSpan) ParseVariable(ctx context.Context) []kv {
	kvs := make([]kv, len(h.kvs))
	copy(kvs, h.kvs)

	if methodName, _ := variable.GetString(ctx, keys.VarMethod); len(methodName) != 0 {
		kvs = append(kvs, kv{"rpc.method", methodName})
	}
	if direction, _ := variable.GetString(ctx, keys.VarDirection); len(direction) != 0 {
		kvs = append(kvs, kv{"hijack", direction})
	}

	if appName, _ := variable.GetString(ctx, keys.VarGovernTargetApp); len(appName) != 0 {
		kvs = append(kvs, kv{"target.app", appName})
	}
	if service, _ := variable.GetString(ctx, keys.VarGovernSourceApp); len(service) != 0 {
		kvs = append(kvs, kv{"caller.app", service})
	}
	dataId, _ := variable.GetString(ctx, keys.VarGovernService)
	kvs = append(kvs, kv{"rpc.service", dataId})
	if len(h.operationName) == 0 {
		h.operationName = dataId
	}

	dp, _ := variable.GetString(ctx, keys.VarDownStreamProtocol)
	if len(dp) != 0 {
		kvs = append(kvs, kv{"downstream.protocol", string(dp)})
	}
	up, _ := variable.GetString(ctx, keys.VarUpStreamProtocol)
	if len(up) != 0 {
		kvs = append(kvs, kv{"upstream.protocol", string(up)})
	} else {
		kvs = append(kvs, kv{"upstream.protocol", string(dp)})
	}
	return kvs
}
