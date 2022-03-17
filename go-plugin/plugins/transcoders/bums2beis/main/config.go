package main

import (
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
