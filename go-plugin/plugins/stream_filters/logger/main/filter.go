package main

import (
	"context"
	"strconv"
	"time"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/protocol/http"
)

type LoggerFilter struct {
	kind     string
	cfg      map[string]string
	rhandler api.StreamReceiverFilterHandler
	shandler api.StreamSenderFilterHandler
	logger   *logger
	tags     [TAG_END]string
}

func NewLoggerFilter(ctx context.Context, cfg map[string]string) *LoggerFilter {
	log := egressLog
	kind := "client"
	if value, ok := config.GetListenerType(ctx); ok && value == "ingress" {
		log = ingressLog
		kind = "server"
	}
	return &LoggerFilter{
		cfg:    cfg,
		logger: log,
		kind:   kind,
	}
}

func (f *LoggerFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.serviceName(headers)
	f.protocolTags(headers)
	return api.StreamFilterContinue
}

func (f *LoggerFilter) Append(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.tags[SPAN_TYPE] = f.kind
	f.tags[APP_NAME] = appName
	span, ok := config.GetSpan(ctx)
	if ok {
		f.tags[TRACEID] = span.TraceId()
		f.tags[SPANID] = span.SpanId()
	}

	requestInfo := f.shandler.RequestInfo()
	startTime := f.shandler.RequestInfo().StartTime()
	endTime := time.Now()
	streamDurationNs := endTime.Sub(startTime).Milliseconds()
	responseReceivedNs := requestInfo.ResponseReceivedDuration().Milliseconds()
	requestReceivedNs := requestInfo.RequestReceivedDuration().Milliseconds()

	processTime := requestReceivedNs // if no response, ignore the network
	if responseReceivedNs > 0 {
		processTime = requestReceivedNs + (streamDurationNs - responseReceivedNs)
	}
	if processTime == 0 {
		processTime = streamDurationNs
	}

	f.tags[START_TIME] = startTime.Format(glogTmFmtWithMS)
	f.tags[END_TIME] = endTime.Format(glogTmFmtWithMS)
	f.tags[DOWNSTEAM_HOST_ADDRESS] = requestInfo.DownstreamRemoteAddress().String()
	f.tags[LISTENER_ADDRESS] = requestInfo.DownstreamLocalAddress().String()
	f.tags[DOWN_PROTOCOL] = string(requestInfo.Protocol())
	f.tags[DURATION] = strconv.FormatInt(streamDurationNs, 10)
	f.tags[MOSN_DURATION] = strconv.FormatInt(processTime, 10)
	f.tags[MOSN_REQUEST_DURATION] = strconv.FormatInt(requestReceivedNs, 10)
	f.tags[MOSN_RSPONSE_DURATION] = strconv.FormatInt(responseReceivedNs, 10)
	if entry := requestInfo.RouteEntry(); entry != nil {
		f.tags[UP_PROTOCOL] = entry.UpstreamProtocol()
	}
	f.tags[UPSTREAM_HOST_ADDRESS] = requestInfo.UpstreamLocalAddress()
	f.tags[RESULT_STATUS] = strconv.Itoa(requestInfo.ResponseCode())
	f.logger.Print(f.tags)
	return api.StreamFilterContinue
}

func (f *LoggerFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.rhandler = handler
}

func (f *LoggerFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.shandler = handler
}

func (f *LoggerFilter) OnDestroy() {
	closeOnce.Do(func() {
		egressLog.Close()
		ingressLog.Close()
	})
}

func (f *LoggerFilter) serviceName(headers api.HeaderMap) {
	if value, ok := headers.Get("X-TARGET-APP"); ok {
		f.tags[SERVICE_NAME] = value
	}
	if value, ok := headers.Get("service"); ok {
		f.tags[SERVICE_NAME] = value
	}
}

func (f *LoggerFilter) protocolTags(headers api.HeaderMap) {
	if hh, ok := headers.(http.RequestHeader); ok {
		f.tags[METHOD_NAME] = string(hh.Method())
		f.tags[REQUEST_URL] = string(hh.RequestURI())
	}

	if value, ok := headers.Get("method"); ok {
		f.tags[METHOD_NAME] = value
	}
	if value, ok := headers.Get("path"); ok {
		f.tags[REQUEST_URL] = value
	}
}
