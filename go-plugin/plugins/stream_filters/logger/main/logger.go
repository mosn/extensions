package main

import (
	"strconv"
	"sync"
)

const (
	glogTmFmtWithMS = "2006-01-02 15:04:05.000"
	gIngressPath    = "/home/admin/logs/mosn_server.log"
	gEgressPath     = "/home/admin/logs/mosn_client.log"
	defaultAppName  = "demo"
)

var (
	traceOnce  sync.Once
	closeOnce  sync.Once
	ingressLog *logger
	egressLog  *logger
	appName    string
)

var (
	maxSize = 1
	maxAge  = 3
)

func initLog(config map[string]string) (err error) {
	// init env
	appName = configValue("app_name", config, defaultAppName)
	// init logger
	ingressLog, err = NewLogger(configValue("ingress_path", config, gIngressPath), config)
	if err != nil {
		return err
	}
	egressLog, err = NewLogger(configValue("egress_path", config, gEgressPath), config)
	if err != nil {
		return err
	}
	return nil
}

func configValue(key string, config map[string]string, defaultValue string) string {
	if v, ok := config[key]; ok {
		return v
	}
	return defaultValue
}

func configIntValue(key string, config map[string]string, defaultValue int) int {
	if v, ok := config[key]; ok {
		if num, err := strconv.Atoi(v); err != nil {
			return num
		}
		return defaultValue
	}
	return defaultValue
}
