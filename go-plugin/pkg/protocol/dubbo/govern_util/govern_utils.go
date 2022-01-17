package govern_util

import (
	"context"
	"mosn.io/api"
)

var (
	ErrorResponseFlag = []api.ResponseFlag{
		api.NoHealthyUpstream,
		api.UpstreamRequestTimeout,
		api.UpstreamLocalReset,
		api.UpstreamRemoteReset,
		api.UpstreamConnectionFailure,
		api.UpstreamConnectionTermination,
		api.NoRouteFound,
		api.DelayInjected,
		api.RateLimited,
		api.ReqEntityTooLarge,
	}

	errorResponseCode = []int{
		api.CodecExceptionCode,
		// UnknownCode
		api.DeserialExceptionCode,
		// SuccessCode
		api.PermissionDeniedCode,
		api.RouterUnavailableCode,
		api.NoHealthUpstreamCode,
		api.UpstreamOverFlowCode,
		api.TimeoutExceptionCode,
		api.LimitExceededCode,
	}
)

var (
	sofaHeadResponseErrorKey = "sofa_head_response_error"
)

var MockDubboMethod string

func AddGovernValue(context context.Context, headers api.HeaderMap, key string, value string) {
	headers.Set(key, value)
}
