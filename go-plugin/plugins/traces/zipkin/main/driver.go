package main

import (
	"encoding/json"
	"fmt"
	"mosn.io/extensions/go-plugin/pkg/common"
	"net"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"

	"mosn.io/pkg/log"
)

func GetTracer(config map[string]interface{}) (*zipkin.Tracer, error) {
	if tracerProvider != nil {
		return tracerProvider, nil
	}
	cfg, err := parseZipkinConfig(config)
	if err != nil {
		return nil, err
	}
	reporterBuilder, ok := GetReportBuilder(cfg.Reporter)
	if !ok {
		return nil, fmt.Errorf("unsupport report type: %s", cfg.Reporter)
	}
	reporter, err := reporterBuilder(cfg)
	if err != nil {
		return nil, err
	}

	sampler, err := zipkin.NewCountingSampler(cfg.SampleRate)
	if err != nil {
		return nil, err
	}

	localIP := common.IpV4
	localEndpoint := &model.Endpoint{
		ServiceName: cfg.ServiceName,
		IPv4:        net.ParseIP(localIP),
		// Port:        localPort,
	}

	tracerProvider, err = zipkin.NewTracer(reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithTraceID128Bit(true),
		zipkin.WithSharedSpans(false),
		zipkin.WithLocalEndpoint(localEndpoint))
	return tracerProvider, err
}

// parseZipkinConfig parse and verify zipkin config
func parseZipkinConfig(config map[string]interface{}) (cfg ZipkinTraceConfig, err error) {
	data, err := json.Marshal(config)
	if err != nil {
		return
	}
	log.DefaultLogger.Infof("[zipkin] [tracer] tracer config: %v", string(data))

	cfg.ServiceName = "mosn"
	if err = json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, cfg.ValidateZipkinConfig()
}
