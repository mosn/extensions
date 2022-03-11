package main

import (
	"context"
	"encoding/json"
	"mosn.io/pkg/protocol/http"

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
	return &BumsFilterFactory{
		config: m,
	}, nil
}

// An implementation of api.StreamFilterChainFactory
type BumsFilterFactory struct {
	config map[string]string
}

func (f *BumsFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewBumsFilter(ctx, f.config)
	// ReceiverFilter, run the filter when receive a request from downstream
	// The FilterPhase can be BeforeRoute or AfterRoute, we use BeforeRoute in this demo
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	// SenderFilter, run the filter when receive a response from upstream
	// In the demo, we are not implement this filter type
	// callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type BumsFilter struct {
	config  map[string]string
	handler api.StreamReceiverFilterHandler
}

// NewBumssFilter returns a BumsFilter, the BumsFilter is an implementation of api.StreamReceiverFilter
// A Filter can implement both api.StreamReceiverFilter and api.StreamSenderFilter.
func NewBumsFilter(ctx context.Context, config map[string]string) *BumsFilter {
	return &BumsFilter{
		config: config,
	}
}

func (f *BumsFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	passed := true
	serviceId := ""
	if _, ok := headers.(http.RequestHeader); ok {
		if buf == nil {
			passed = false
		} else {
			var body map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &body)
			if err != nil {
				passed = false
				log.DefaultLogger.Errorf("Unmarshal ERR %s", err)
			} else {
				if body["head"] != nil {
					_v, _ := json.Marshal(body["head"])
					var bodyHead map[string]string
					json.Unmarshal(_v, &bodyHead)
					tranCode := bodyHead["tranCode"]
					serviceCode := bodyHead["serviceCode"]
					serviceScene := bodyHead["serviceScene"]

					serviceId = getServiceId(tranCode, serviceCode, serviceScene)
					if serviceId == "" {
						passed = false
						log.DefaultLogger.Errorf("[stream_filter][bums] Not Found ServiceId")
					}
				} else {
					passed = false
				}
			}

		}
	}

	if !passed {
		return api.StreamFilterStop
	}
	// inject http header
	headers.Set("X-TARGET-APP", serviceId)
	headers.Set("X-SERVICE-TYPE", "springcloud")

	return api.StreamFilterContinue
}

func (f *BumsFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.handler = handler
}

func (f *BumsFilter) OnDestroy() {}

func getServiceId(tranCode string, serviceCode string, serviceScene string) string {
	return ""
}
