package main

import (
	"context"
	"encoding/json"

	"mosn.io/api"
)

func CreateFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	b, _ := json.Marshal(conf)
	m := make(map[string]string)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	var err error
	traceOnce.Do(func() {
		err = initLog(m)
	})
	if err != nil {
		return nil, err
	}
	return &LoggerFilterFactory{
		config: m,
	}, nil
}

type LoggerFilterFactory struct {
	config map[string]string
}

func (f *LoggerFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewLoggerFilter(ctx, f.config)
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}
