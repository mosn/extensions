package trace

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"mosn.io/api"
	"mosn.io/pkg/log"
)

const (
	parentSpanID = "-1"
)

var MosnProcessFailedFlags = api.NoHealthyUpstream | api.NoRouteFound | api.UpstreamLocalReset |
	api.FaultInjected | api.RateLimited | api.DownStreamTerminate | api.ReqEntityTooLarge

var _ api.Span = NoopSpan{}

type NoopSpan struct {
	tid string
	sid string
}

func NewNooSpan() NoopSpan {
	return NoopSpan{
		tid: uuid.NewV4().String(),
		sid: uuid.NewV4().String(),
	}
}

func (n NoopSpan) TraceId() string {
	return n.tid
}

func (n NoopSpan) SpanId() string {
	return n.sid
}

func (NoopSpan) ParentSpanId() string {
	return parentSpanID
}

func (NoopSpan) FinishSpan() {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported FinishSpan")
	}
}

func (NoopSpan) SetOperation(operation string) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported SetOperation [%s]", operation)
	}
}

func (NoopSpan) SetTag(key uint64, value string) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported SetTag [%d]-[%s]", key, value)
	}
}

func (NoopSpan) Tag(key uint64) string {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported Tag [%s]", key)
	}
	return ""
}

func (NoopSpan) InjectContext(requestHeaders api.HeaderMap, requestInfo api.RequestInfo) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported InjectContext")
	}
}

func (NoopSpan) SetRequestInfo(requestInfo api.RequestInfo) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported SetRequestInfo")
	}
}

func (NoopSpan) SpawnChild(operationName string, _ time.Time) api.Span {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[Noop] [tracer] [span] Unsupported SpawnChild [%s]", operationName)
	}
	return nil
}
