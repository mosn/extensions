package dubbo

import (
	"context"
	"errors"
	"net/http"

	hessian "github.com/apache/dubbo-go-hessian2"
	"mosn.io/api"
)

type StatusMapping struct{}

func (m StatusMapping) MappingHeaderStatusCode(ctx context.Context, headers api.HeaderMap) (int, error) {
	cmd, ok := headers.(api.XRespFrame)
	if !ok {
		return 0, errors.New("no response status in headers")
	}
	code := cmd.GetStatusCode()
	// TODO: more accurate mapping
	switch byte(code) {
	case hessian.Response_OK:
		return http.StatusOK, nil
	default:
		return http.StatusInternalServerError, nil
	}
}
