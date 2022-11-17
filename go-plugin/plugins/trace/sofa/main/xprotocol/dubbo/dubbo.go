package dubbo

import (
	"bytes"
	"context"
	"strings"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/extensions/go-plugin/plugins/trace/sofa/main/generator"
)

const (
	dubboSofaTracerItemSep   = "&"
	dubboSofaTracerItemKVSep = "="

	HeaderDubboSofaTracerSpanId                   = "spid"
	HeaderDubboSofaTracerTraceId                  = "tcid"
	HeaderDubboSofaTracer                         = "dubbo.rpc.sofa.tracer"
	MethodKey                                     = "method"
	ProtocolName                 api.ProtocolName = "dubbo" // protocol
)

func init() {
	generator.RegisterDelegate(ProtocolName, DubboDelegate)
}

func DubboDelegate(ctx context.Context, frame api.XFrame, span api.Span) {
	request := frame.GetHeader()
	lType, _ := config.GetListenerType(ctx)
	var traceId string
	var spanId string
	dubboSofaTracer, found := request.Get(HeaderDubboSofaTracer)
	if found {
		items := strings.Split(dubboSofaTracer, dubboSofaTracerItemSep)
		for _, item := range items {
			kv := strings.Split(item, dubboSofaTracerItemKVSep)
			if len(kv) != 2 {
				continue
			}
			if kv[0] == HeaderDubboSofaTracerTraceId {
				traceId = kv[1]
			} else if kv[0] == HeaderDubboSofaTracerSpanId {
				spanId = kv[1]
			}
		}
	} else {
		traceId, _ = request.Get(generator.TRACER_ID_KEY)
		spanId, _ = request.Get(generator.RPC_ID_KEY)
	}
	if len(traceId) == 0 {
		span.SetTag(generator.SPAN_ID, "0")
		span.SetTag(generator.TRACE_ID, generator.IdGen().GenerateTraceId())
	} else {
		span.SetTag(generator.TRACE_ID, traceId)
		if lType == "ingress" {
			generator.AddSpanIdGenerator(generator.NewSpanIdGenerator(traceId, spanId))
		} else {
			span.SetTag(generator.PARENT_SPAN_ID, spanId)
			spanKey := &generator.SpanKey{TraceId: traceId, SpanId: spanId}
			if spanIdGenerator := generator.GetSpanIdGenerator(spanKey); spanIdGenerator != nil {
				spanId = spanIdGenerator.GenerateNextChildIndex()
			}
		}
		span.SetTag(generator.SPAN_ID, spanId)
	}

	if lType == "EGRESS" {
		if appName, found := request.Get(generator.APP_NAME_KEY); found {
			span.SetTag(generator.CALLER_APP_NAME, appName)
		}
	}
	if methodName, found := request.Get(MethodKey); found {
		span.SetTag(generator.METHOD_NAME, methodName)
	}

	dubboServiceName, _ := request.Get(generator.XPROTOCOL_DUBBO_META_INTERFACE)
	dubboServiceVersion, _ := request.Get(generator.XPROTOCOL_DUBBO_META_VERSION)
	dubboServiceGroup, _ := request.Get(generator.XPROTOCOL_DUBBO_META_GROUP)
	if "0.0.0" == dubboServiceVersion {
		dubboServiceVersion = ""
	}
	dataId := BuildDubboDataId(dubboServiceName, dubboServiceVersion, dubboServiceGroup)
	if dubboServiceName != "" {
		span.SetTag(generator.SERVICE_NAME, dubboServiceName)
	} else if serviceName, found := request.Get(generator.SERVICE_KEY); found {
		span.SetTag(generator.SERVICE_NAME, serviceName)
	}

	if baggageData, found := request.Get(generator.SOFA_TRACE_BAGGAGE_DATA); found {
		span.SetTag(generator.BAGGAGE_DATA, baggageData)
	}

	span.SetTag(generator.APP_SERVICE_NAME, dataId)
}

func BuildDubboDataId(name, version, group string) string {
	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	buffer.WriteString(name)
	if version != "" {
		buffer.WriteString(":")
		buffer.WriteString(version)
	}
	if group != "" {
		buffer.WriteString(":")
		buffer.WriteString(group)
	}
	buffer.WriteString(generator.XPROTOCOL_TYPE_DUBBO_DATAID_SUFFIX)
	return buffer.String()
}
