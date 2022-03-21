package main

import (
	"crypto/md5"
	"encoding/json"

	"mosn.io/extensions/go-plugin/pkg/transcoder/bumscd"
)

var (
	// 单线程安全
	configManager ConfigManager
)

type ConfigManager struct {
	cmap map[string]*bumscd.Config
}

func NewConfig(info string) (*ConfigManager, error) {
	return &ConfigManager{
		cmap: make(map[string]*bumscd.Config),
	}, nil
}

func (cm *ConfigManager) GetLatestRelation(info string) (*bumscd.Config, error) {
	rmd5 := md5.New().Sum(bumscd.S2B(info))
	if cfg, ok := cm.cmap[bumscd.B2S(rmd5)]; ok {
		return cfg, nil
	}

	cfg := &bumscd.Config{}
	if err := json.Unmarshal(bumscd.S2B(info), cfg); err != nil {
		return nil, err
	}
	cm.cmap[bumscd.B2S(rmd5)] = cfg
	return cfg, nil
}
