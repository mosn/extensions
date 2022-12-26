package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	language_agent "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/pkg/trace"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

type kv struct {
	Key   string
	Value string
}

type SkySpan struct {
	trace.NoopSpan
	tid       string
	sid       string
	psid      string
	startTime time.Time
	ctx       context.Context

	operationName string
	kvs           []kv
	provider      *go2sky.Tracer
	entrySpan     go2sky.Span
	entryctx      context.Context
	exitSpan      go2sky.Span
}

func NewSpan(ctx context.Context, startTime time.Time, provider *go2sky.Tracer) *SkySpan {
	h := &SkySpan{
		startTime: startTime,
		ctx:       ctx,
		provider:  provider,
		kvs:       make([]kv, 0, 10),
	}
	return h
}

func (h *SkySpan) TraceId() string {
	return h.tid
}

func (h *SkySpan) SpanId() string {
	return h.sid
}

func (h *SkySpan) ParentSpanId() string {
	return h.psid
}

func (h *SkySpan) InjectContext(headers api.HeaderMap, reqInfo api.RequestInfo) {
	upstreamLocalAddress := reqInfo.UpstreamLocalAddress()
	if header, ok := headers.(http.RequestHeader); ok {
		requestURI := string(header.RequestURI())
		url := strings.Join([]string{"http://", upstreamLocalAddress, requestURI}, "")
		h.kvs = append(h.kvs, kv{string(go2sky.TagURL), url})
		h.kvs = append(h.kvs, kv{string(go2sky.TagHTTPMethod), string(header.Method())})
		h.operationName = requestURI
		exit, err := h.provider.CreateExitSpan(h.entryctx, requestURI, upstreamLocalAddress, func(Value string) error {
			headers.Set(propagation.Header, Value)
			return nil
		})
		if err != nil {
			log.DefaultLogger.Errorf("[SkyWalking] [tracer] [http1] create exit span error, err: %v", err)
			return
		}
		exit.SetComponent(MOSNComponentID)
		exit.SetSpanLayer(language_agent.SpanLayer_Http)
		h.exitSpan = exit
	} else {
		exit, err := h.provider.CreateExitSpan(h.entryctx, "mosn", upstreamLocalAddress, func(Value string) error {
			headers.Set(propagation.Header, Value)
			return nil
		})
		if err != nil {
			log.DefaultLogger.Errorf("[SkyWalking] [tracer] [http1] create exit span error, err: %v", err)
			return
		}
		exit.SetComponent(MOSNComponentID)
		exit.SetSpanLayer(language_agent.SpanLayer_RPCFramework)
		h.exitSpan = exit
	}
}

func (h *SkySpan) SetRequestInfo(reqInfo api.RequestInfo) {
	h.setRequestInfo(reqInfo)
	responseCode := strconv.Itoa(reqInfo.ResponseCode())

	// end exit span (upstream)
	if h.exitSpan != nil {
		exit := h.exitSpan
		if reqInfo.ResponseCode() != api.SuccessCode {
			ok := reqInfo.GetResponseFlag(trace.MosnProcessFailedFlags)
			exit.Tag(go2sky.TagStatusCode, strconv.Itoa(reqInfo.ResponseCode()))
			if ok {
				exit.Error(time.Now(), "mosn error")
			} else {
				exit.Error(time.Now(), "biz error")
			}
		} else {
			exit.Tag(go2sky.TagStatusCode, responseCode)
		}
		kvs := h.ParseVariable(h.ctx)
		h.addTag(exit, kvs)
		exit.SetOperationName(h.operationName)
		exit.End()
		h.log(kvs, go2sky.SpanTypeExit)
	}

	// entry span (downstream)
	entry := h.entrySpan
	if reqInfo.ResponseCode() != api.SuccessCode {
		ok := reqInfo.GetResponseFlag(trace.MosnProcessFailedFlags)
		entry.Tag(go2sky.TagStatusCode, strconv.Itoa(reqInfo.ResponseCode()))
		if ok {
			h.entrySpan.Error(time.Now(), "mosn error")
		} else {
			h.entrySpan.Error(time.Now(), "biz error")
		}
	} else {
		entry.Tag(go2sky.TagStatusCode, responseCode)
	}
}

func (h *SkySpan) setRequestInfo(reqInfo api.RequestInfo) {
	if host := reqInfo.UpstreamHost(); host != nil {
		h.kvs = append(h.kvs, kv{"upstream.address", host.AddressString()})
	}
	if addr := reqInfo.DownstreamRemoteAddress(); addr != nil {
		h.kvs = append(h.kvs, kv{"downstream.address", addr.String()})
	}
	h.kvs = append(h.kvs, kv{"request.size", strconv.Itoa(int(reqInfo.BytesReceived()))})
	h.kvs = append(h.kvs, kv{"respone.size", strconv.Itoa(int(reqInfo.BytesSent()))})
	h.kvs = append(h.kvs, kv{"duration", strconv.Itoa(int(reqInfo.Duration().Nanoseconds()))})
	process := reqInfo.ProcessTimeDuration().Nanoseconds()
	if process == 0 {
		process = reqInfo.RequestFinishedDuration().Nanoseconds()
	}
	h.kvs = append(h.kvs, kv{"mosn.process.duration", strconv.Itoa(int(process))})
	h.kvs = append(h.kvs, kv{"mosn.process.request.duration", strconv.Itoa(int(reqInfo.RequestFinishedDuration().Nanoseconds()))})
	h.kvs = append(h.kvs, kv{"mosn.process.respone.duration", strconv.Itoa(int(reqInfo.ResponseReceivedDuration().Nanoseconds()))})
	if reqInfo.ResponseCode() != api.SuccessCode {
		ok := reqInfo.GetResponseFlag(trace.MosnProcessFailedFlags)
		h.kvs = append(h.kvs, kv{"mosn.process.fail", strconv.FormatBool(ok)})
	}
}

func (h *SkySpan) FinishSpan() {
	if h.entrySpan != nil {
		kvs := h.ParseVariable(h.ctx)
		h.addTag(h.entrySpan, kvs)
		h.entrySpan.SetOperationName(h.operationName)
		currentIP, _ := trace.GetOutboundIP()
		h.entrySpan.SetPeer(currentIP)
		h.entrySpan.End()
		h.log(kvs, go2sky.SpanTypeEntry)
	}
}

func (h *SkySpan) addTag(span go2sky.Span, kvs []kv) {
	for _, kv := range kvs {
		span.Tag(go2sky.Tag(kv.Key), kv.Value)
	}
}

func (h *SkySpan) SetOperation(operation string) {
	h.operationName = operation
}

func (h *SkySpan) log(kvs []kv, stype go2sky.SpanType) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		kvs, _ := json.Marshal(h.kvs)
		log.DefaultLogger.Debugf("trace:%s pid:%s parentid:%s operationName:%s,kvs:%s", h.tid, h.sid, h.psid, h.operationName, kvs)
	}
}

func (h *SkySpan) CreateLocalHttpSpan(ctx context.Context, header http.RequestHeader, entry go2sky.Span) {
	h.entryctx = ctx
	h.entrySpan = entry
	// TODO parent span
	// h.psid = span.SpanContext().SpanID().String()
	// h.sid = strconv.Itoa(int(go2sky.SpanID(ctx)))
	h.tid = go2sky.TraceID(ctx)
	requestURI := string(header.RequestURI())
	url := strings.Join([]string{"http://", string(header.Host()), string(header.RequestURI())}, "")
	h.kvs = append(h.kvs, kv{"caller.url", url})
	h.kvs = append(h.kvs, kv{"caller.method", string(header.Method())})
	h.operationName = requestURI
}

func (h *SkySpan) CreateLocalRpcSpan(ctx context.Context, entry go2sky.Span) {
	h.entryctx = ctx
	h.entrySpan = entry
	// TODO parent span
	// h.psid = span.SpanContext().SpanID().String()
	// h.sid = strconv.Itoa(int(go2sky.SpanID(ctx)))
	h.tid = go2sky.TraceID(ctx)
}

func (h *SkySpan) ParseVariable(ctx context.Context) []kv {
	kvs := make([]kv, len(h.kvs))
	copy(kvs, h.kvs)
	/*
			if methodName, _ := variable.GetString(ctx, govern.VarGovernMethod); len(methodName) != 0 {
				kvs = append(kvs, kv{"rpc.method", methodName})
			}
			if direction, _ := variable.GetString(ctx, govern.VarGovernDirection); len(direction) != 0 {
				kvs = append(kvs, kv{"hijack", direction})
			}

			if appName, _ := variable.GetString(ctx, govern.VarGovernTargetApp); len(appName) != 0 {
				kvs = append(kvs, kv{"target.app", appName})
			}
			if service, _ := variable.GetString(ctx, govern.VarGovernSourceApp); len(service) != 0 {
				kvs = append(kvs, kv{"caller.app", service})
			}
		dataId, _ := variable.GetString(ctx, govern.VarGovernService)
		if len(h.operationName) == 0 {
			h.operationName = dataId
		}
	*/
	dp, _ := config.GetDownstreamProtocol(h.ctx)
	if len(dp) != 0 {
		kvs = append(kvs, kv{"downstream.protocol", string(dp)})
	}
	up, _ := config.GetUpstreamProtocol(h.ctx)
	if len(up) != 0 {
		kvs = append(kvs, kv{"upstream.protocol", string(up)})
	} else {
		kvs = append(kvs, kv{"upstream.protocol", string(dp)})
	}

	return kvs
}
