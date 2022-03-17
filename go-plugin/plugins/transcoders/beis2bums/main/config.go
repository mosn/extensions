package main

import (
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
