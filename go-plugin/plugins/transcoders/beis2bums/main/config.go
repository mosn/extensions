package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"mosn.io/extensions/go-plugin/pkg/transcoder/bumsbeis"
)

type Bums2BeisConfig struct {
	UniqueId    string                    `json:"uniqueId"`
	Path        string                    `json:"path"`
	Method      string                    `json:"method"`
	GWName      string                    `json:"gw"`
	ReqMapping  *bumsbeis.Beis2BumsConfig `json:"-"`
	RespMapping *bumsbeis.Bums2BeisConfig `json:"resp_mapping"`
}

var (
	// 单线程安全
	configManager = NewConfig()
)

type ConfigManager struct {
	cmap map[string][]*Bums2BeisConfig
}

func NewConfig() *ConfigManager {
	return &ConfigManager{
		cmap: make(map[string][]*Bums2BeisConfig),
	}
}

func (cm *ConfigManager) GetLatestRelation(info string) (*Bums2BeisConfig, error) {
	rmd5 := md5.New().Sum([]byte(info))
	if cfg, ok := cm.cmap[string(rmd5)]; ok {
		return cfg[0], nil
	}

	var cfg []*Bums2BeisConfig
	if err := json.Unmarshal([]byte(info), &cfg); err != nil {
		return nil, err
	}
	if len(cfg) != 1 {
		return nil, fmt.Errorf("the length of configs is illage")
	}
	cfg[0].ReqMapping = &bumsbeis.Beis2BumsConfig{
		Path:   cfg[0].Path,
		Method: cfg[0].Method,
		GWName: cfg[0].GWName,
	}
	cm.cmap[string(rmd5)] = cfg
	return cfg[0], nil
}
