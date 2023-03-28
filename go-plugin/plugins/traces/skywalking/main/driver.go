package main

import (
	"encoding/json"
	"errors"
	"mosn.io/extensions/go-plugin/pkg/common"
	"strconv"
	"sync"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"google.golang.org/grpc/credentials"
)

var (
	ReporterCfgErr       = errors.New("SkyWalking tracer support only log and gRPC reporter")
	BackendServiceCfgErr = errors.New("SkyWalking tracer must configure the backend_service")
	tracerProvider       *go2sky.Tracer
	skyLock              sync.RWMutex
	SkyDriverName        = "SkyWalking"
)

func GetTracer(config map[string]interface{}) (*go2sky.Tracer, error) {
	skyLock.Lock()
	defer skyLock.Unlock()
	if tracerProvider != nil {
		return tracerProvider, nil
	}
	cfg, err := parseAndVerifySkyTracerConfig(config)
	if err != nil {
		return nil, err
	}

	var r go2sky.Reporter
	if cfg.Reporter == LogReporter {
		r, err = reporter.NewLogReporter()
		if err != nil {
			return nil, err
		}
	} else if cfg.Reporter == GRPCReporter {
		// opts
		var opts []reporter.GRPCReporterOption
		opts = append(opts, reporter.WithLogger(NewDefaultLogger()))
		// max send queue size
		if size, _ := strconv.Atoi(cfg.MaxSendQueueSize); size > 0 {
			opts = append(opts, reporter.WithMaxSendQueueSize(size))
		}
		// auth
		if cfg.Authentication != "" {
			opts = append(opts, reporter.WithAuthentication(cfg.Authentication))
		}
		// tls
		if cfg.TLS.CertFile != "" {
			cReds, err := credentials.NewClientTLSFromFile(cfg.TLS.CertFile, cfg.TLS.ServerNameOverride)
			if err != nil {
				return nil, err
			}
			opts = append(opts, reporter.WithTransportCredentials(cReds))
		}
		r, err = reporter.NewGRPCReporter(cfg.BackendService, opts...)
		if err != nil {
			return nil, err
		}
	}

	serviceName := cfg.ServiceName
	if len(serviceName) > 4 {
		serviceName = serviceName[0:4] + "::" + serviceName
	} else {
		serviceName = serviceName + "::" + serviceName
	}
	currentIP := common.IpV4
	if cfg.VmMode != "" {
		tracerProvider, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(r), go2sky.WithInstance(currentIP))
		return tracerProvider, err
	}
	skyInstanceId := cfg.PodName
	tracerProvider, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(r), go2sky.WithInstance(skyInstanceId))
	return tracerProvider, err
}

func parseAndVerifySkyTracerConfig(cfg map[string]interface{}) (config SkyWalkingTraceConfig, err error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return config, err
	}
	// set default value
	config.Reporter = LogReporter
	config.ServiceName = DefaultServiceName
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	if config.Reporter != LogReporter && config.Reporter != GRPCReporter {
		return config, ReporterCfgErr
	}

	if config.Reporter == GRPCReporter && config.BackendService == "" {
		return config, BackendServiceCfgErr
	}
	return config, nil
}
