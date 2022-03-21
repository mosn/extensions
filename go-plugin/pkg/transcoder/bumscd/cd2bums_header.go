package bumscd

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

func (cd2bm *Cd2Bums) HeadRespone() (api.HeaderMap, error) {
	respHeader := &fasthttp.ResponseHeader{}
	respHeader.Set("Content-Type", "application/json")
	cd2bm.header.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			respHeader.Set(key, value)
		}
		return true
	})

	if code, ok := cd2bm.header.Get("x-mosn-status"); ok {
		statusCode, err := strconv.Atoi(code)
		if err == nil {
			respHeader.SetStatusCode(statusCode)
		} else {
			log.DefaultContextLogger.Warnf(cd2bm.ctx, "the atoi of statuscode failed. err:%s", err)
		}
	}

	// beis数据解析
	return http.ResponseHeader{respHeader}, nil
}

func (cd2bm *Cd2Bums) HeadRequest() (api.HeaderMap, error) {
	reqHeader := &fasthttp.RequestHeader{}
	reqHeader.Set("Content-Type", "application/json")
	cd2bm.header.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			reqHeader.Set(key, value)
		}
		return true
	})
	reqHeader.Set("x-mosn-method", cd2bm.config.Method)
	reqHeader.Set("x-mosn-path", cd2bm.config.Path)
	reqHeader.Set("X-TARGET-APP", cd2bm.config.GWName)

	return http.RequestHeader{reqHeader}, nil
}
