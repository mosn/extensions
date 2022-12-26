package main

import (
	"log"
	"strconv"
	"time"

	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/openzipkin/zipkin-go/reporter/kafka"
)

var (
	factory = make(map[string]ReporterBuilder)
)

func init() {
	factory[ZipkinHttpReport] = HttpReporterBuilder
	factory[ZipkinKafkaReport] = KafkaReporterBuilder
}

func GetReportBuilder(typ string) (ReporterBuilder, bool) {
	if v, ok := factory[typ]; ok {
		return v, ok
	}
	return nil, false
}

type ReporterBuilder func(ZipkinTraceConfig) (reporter.Reporter, error)

func HttpReporterBuilder(cfg ZipkinTraceConfig) (reporter.Reporter, error) {
	opts := make([]http.ReporterOption, 0, 4)
	opts = append(opts, http.Logger(log.Default()))
	if cfg.HttpConfig.Timeout > 0 {
		opts = append(opts, http.Timeout(time.Second*time.Duration(cfg.HttpConfig.Timeout)))

	}
	if cfg.HttpConfig.BatchInterval > 0 {
		opts = append(opts, http.BatchInterval(time.Second*time.Duration(cfg.HttpConfig.BatchInterval)))
	}
	if size, _ := strconv.Atoi(cfg.HttpConfig.BatchSize); size > 0 {
		opts = append(opts, http.BatchSize(size))
	}
	return http.NewReporter(cfg.HttpConfig.Address, opts...), nil
}

func KafkaReporterBuilder(cfg ZipkinTraceConfig) (reporter.Reporter, error) {
	opts := make([]kafka.ReporterOption, 0, 2)
	opts = append(opts, kafka.Logger(log.Default()))
	if cfg.KafkaConfig.Topic != "" {
		opts = append(opts, kafka.Topic(cfg.KafkaConfig.Topic))
	}
	return kafka.NewReporter(cfg.KafkaConfig.Addresses, opts...)
}
