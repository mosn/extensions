package bumsbeis

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

func (bibm *Beis2Bums) HeadRespone() (api.HeaderMap, error) {
	respHeader := &fasthttp.ResponseHeader{}
	respHeader.Set("Content-Type", "application/json")
	bibm.header.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			respHeader.Set(key, value)
		}
		return true
	})

	if code, ok := bibm.header.Get("x-mosn-status"); ok {
		statusCode, err := strconv.Atoi(code)
		if err == nil {
			respHeader.SetStatusCode(statusCode)
		} else {
			log.DefaultContextLogger.Warnf(bibm.ctx, "the atoi of statuscode failed. err:%s", err)
		}
	}

	// beis数据解析
	br := bibm.header.(*beis.Response)
	respHeader.Set("VersionId", br.VersionID)
	respHeader.Set("OrigSender", br.OrigSender)
	respHeader.Set("CtrlBits", br.CtrlBits)
	respHeader.Set("AreaCode", br.AreaCode)
	return http.ResponseHeader{respHeader}, nil
}

func (bibm *Beis2Bums) HeadRequest() (api.HeaderMap, error) {
	reqHeader := &fasthttp.RequestHeader{}
	reqHeader.Set("Content-Type", "application/json")
	bibm.header.Range(func(key, value string) bool {
		if key != "Content-Length" && key != "Accept:" {
			reqHeader.Set(key, value)
		}
		return true
	})
	reqHeader.Set("x-mosn-method", bibm.config.Method)
	reqHeader.Set("x-mosn-path", bibm.config.Path)
	reqHeader.Set("X-TARGET-APP", bibm.config.GWName)

	// beis数据解析
	br := bibm.header.(*beis.Request)
	reqHeader.Set("VersionId", br.VersionID)
	reqHeader.Set("OrigSender", br.OrigSender)
	reqHeader.Set("CtrlBits", br.CtrlBits)
	reqHeader.Set("AreaCode", br.AreaCode)
	return http.RequestHeader{reqHeader}, nil
}
