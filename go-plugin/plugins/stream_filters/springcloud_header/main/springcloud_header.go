package main

import (
	"context"
	"encoding/json"

	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
)

// define a function named: CreateFilterFactory, do not need init to register
func CreateFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	b, _ := json.Marshal(conf)
	m := make(map[string]string)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &SpringCloudHeaderFilterFactory{
		config: m,
	}, nil
}

// An implementation of api.StreamFilterChainFactory
type SpringCloudHeaderFilterFactory struct {
	config map[string]string
}

func (f *SpringCloudHeaderFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewSpringCloudHeadersFilter(ctx, f.config)
	// ReceiverFilter, run the filter when receive a request from downstream
	// The FilterPhase can be BeforeRoute or AfterRoute, we use BeforeRoute in this demo
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	// SenderFilter, run the filter when receive a response from upstream
	// In the demo, we are not implement this filter type
	// callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type SpringCloudHeaderFilter struct {
	config  map[string]string
	handler api.StreamReceiverFilterHandler
}

// NewSpringCloudHeadersFilter returns a SpringCloudHeaderFilter, the SpringCloudHeaderFilter is an implementation of api.StreamReceiverFilter
// A Filter can implement both api.StreamReceiverFilter and api.StreamSenderFilter.
func NewSpringCloudHeadersFilter(ctx context.Context, config map[string]string) *SpringCloudHeaderFilter {
	return &SpringCloudHeaderFilter{
		config: config,
	}
}

type SpringCloudHeader struct {
	TargetApp   string `json:"targetApp"`
	ServiceType string `json:"serviceType"`
}

func (f *SpringCloudHeaderFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	passed := true
	// check headers contains `X-TARGET-APP` already ?
	if app, ok := headers.Get("X-TARGET-APP"); !ok || app == "" {
		passed = false
	}
	// check headers contains `X-SERVICE-TYPE` already ?
	if serviceType, ok := headers.Get("X-SERVICE-TYPE"); !ok || serviceType == "" {
		passed = false
	}

	if passed {
		// headers already contains `X-TARGET-APP` and `X-SERVICE-TYPE`
		return api.StreamFilterContinue
	}

	// try to decode body
	var request SpringCloudHeader
	if err := json.Unmarshal(buf.Bytes(), &request); err != nil {
		log.DefaultContextLogger.Warnf(ctx, "[streamfilter][springcloud_header] Unmarshal body stream to spring cloud header failed, content %s", buf.Bytes())
		f.handler.SendHijackReply(403, headers)
		return api.StreamFilterStop
	}

	if request.TargetApp == "" {
		log.DefaultContextLogger.Warnf(ctx, "[streamfilter][springcloud_header] targetApp is required, content %s", buf.Bytes())
		f.handler.SendHijackReply(403, headers)
		return api.StreamFilterStop
	}

	if request.ServiceType == "" {
		request.ServiceType = "springcloud"
	}

	// inject http header
	headers.Set("X-TARGET-APP", request.TargetApp)
	headers.Set("X-SERVICE-TYPE", request.ServiceType)

	return api.StreamFilterContinue
}

func (f *SpringCloudHeaderFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.handler = handler
}

func (f *SpringCloudHeaderFilter) OnDestroy() {}
