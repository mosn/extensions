package main

import (
	"errors"
	"strings"
)

const (
	ZipkinKafkaReport string = "kafka"
	ZipkinHttpReport  string = "http"
	ZipkinTracer             = "Zipkin"
)

type ZipkinTraceConfig struct {
	ServiceName string             `json:"service_name"`
	Reporter    string             `json:"reporter"`
	SampleRate  float64            `json:"sample_rate"`
	HttpConfig  *HttpReportConfig  `json:"http"`
	KafkaConfig *KafkaReportConfig `json:"kafka"`
}

type HttpReportConfig struct {
	Timeout       int    `json:"timeout"`
	BatchSize     int    `json:"batch_size"`
	BatchInterval int    `json:"batch_interval"`
	Address       string `json:"address"`
}

type KafkaReportConfig struct {
	Topic     string   `json:"topic"`
	Address   string   `json:"address"`
	Addresses []string `json:"addresses"`
}

func (z *ZipkinTraceConfig) ValidateZipkinConfig() error {
	if z.SampleRate > 1 || z.SampleRate < 0 {
		return errors.New("sample rate should between 1.0 and 0.0")
	}

	switch z.Reporter {
	case ZipkinHttpReport:
		if len(z.HttpConfig.Address) == 0 {
			return errors.New("http config only support single address")
		}
	case ZipkinKafkaReport:
		if len(z.KafkaConfig.Address) == 0 {
			return errors.New("kafka config address can't be empty")
		}
		z.KafkaConfig.Addresses = strings.Split(z.KafkaConfig.Address, ",")
	default:
		return errors.New("not support")
	}
	return nil
}
