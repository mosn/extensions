package main

import (
	"context"
	"encoding/json"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/config"
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
	return &DynamicFilterFactory{
		config: m,
	}, nil
}

// An implementation of api.StreamFilterChainFactory
type DynamicFilterFactory struct {
	config map[string]string
}

func (f *DynamicFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewDynamicFilter(ctx, f.config)
	// ReceiverFilter, run the filter when receive a request from downstream
	// The FilterPhase can be BeforeRoute or AfterRoute, we use BeforeRoute in this demo
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	// SenderFilter, run the filter when receive a response from upstream
	// In the demo, we are not implement this filter type
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

// What DynamicFilter do:
// the request will be passed only if the request headers contains key&value matched in the config
type DynamicFilter struct {
	config   map[string]string
	rhandler api.StreamReceiverFilterHandler
	shandler api.StreamSenderFilterHandler
}

// NewDynamicFilter returns a DynamicFilter, the DynamicFilter is an implementation of api.StreamReceiverFilter
// A Filter can implement both api.StreamReceiverFilter and api.StreamSenderFilter.
func NewDynamicFilter(ctx context.Context, config map[string]string) *DynamicFilter {
	return &DynamicFilter{
		config: config,
	}
}

func (f *DynamicFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	conf, ok := config.GlobalExtendConfigByContext(ctx, "config")
	log.DefaultContextLogger.Infof(ctx, "get dynamic conf:%s ok:%v", conf, ok)
	headers.Set("x-request-proxy", conf)
	return api.StreamFilterContinue
}

func (f *DynamicFilter) Append(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	conf, ok := config.GlobalExtendConfigByContext(ctx, "config")
	log.DefaultContextLogger.Infof(ctx, "get dynamic conf:%s ok:%v", conf, ok)
	headers.Set("x-response-proxy", conf)
	return api.StreamFilterContinue
}

func (f *DynamicFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.rhandler = handler
}

func (f *DynamicFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.shandler = handler
}

func (f *DynamicFilter) OnDestroy() {}
