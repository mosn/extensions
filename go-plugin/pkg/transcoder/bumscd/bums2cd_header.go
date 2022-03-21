package bumscd

import (
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/cd"
)

func (bm2cd *Bums2Cd) HeadRequest() (api.HeaderMap, error) {
	headers := &cd.Request{}
	bm2cd.header.Range(func(key, value string) bool {
		headers.Set(key, value)
		return true
	})
	headers.Set("serviceKey", bm2cd.config.GWName)
	return headers, nil
}

func (bm2cd *Bums2Cd) HeadRespone() (api.HeaderMap, error) {
	headers := &cd.Response{}
	bm2cd.header.Range(func(key, value string) bool {
		headers.Set(key, value)
		return true
	})
	headers.Set("serviceKey", bm2cd.config.GWName)
	return headers, nil
}
