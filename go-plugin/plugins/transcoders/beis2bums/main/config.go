package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"mosn.io/extensions/go-plugin/pkg/transcoder/bumsbeis"
)

type Bums2BeisConfig struct {
	UniqueId     string                    `json:"uniqueId"`
	ServiceCode  string                    `json:"service_code"`
	ServiceScene string                    `json:"service_scene"`
	GWName       string                    `json:"gw"`
	ReqMapping   *bumsbeis.Bums2BeisConfig `json:"req_mapping"`
	// RespMapping  *bumsbeis.Bums2BeisConfig `json:"resp_mapping"`
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
	cm.cmap[string(rmd5)] = cfg
	return cfg[0], nil
}
