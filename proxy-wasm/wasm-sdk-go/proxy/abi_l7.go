package proxy

import "github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"

//export proxy_on_request_headers
func proxyOnRequestHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_headers")
	}
	this.setActiveContextID(contextID)
	var header Header
	if numHeaders > 0 {
		hs, err := getHttpRequestHeaders()
		if err != nil {
			log.Errorf("failed to get request headers: %v", err)
			return types.ActionContinue
		}
		header = &CommonHeader{m: hs}
		// update context header
		attr := ctx.(attribute)
		attr.set(types.AttributeKeyHeaderHolder, header)
	}

	if endOfStream {
		return ctx.OnDownStreamReceived(header, NewBuffer(0), nil)
	}

	return types.ActionContinue
}

//export proxy_on_request_body
func proxyOnRequestBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_body")
	}
	this.setActiveContextID(contextID)

	attr := ctx.(attribute)

	var body Buffer
	if bodySize > 0 {
		bodyBytes, err := getHttpRequestBody(0, bodySize)
		if err != nil {
			log.Errorf("failed to get request body: %v", err)
			return types.ActionContinue
		}

		body = WrapBuffer(bodyBytes)
		// update context body buffer
		attr.set(types.AttributeKeyBufferHolder, body)
	}

	if endOfStream {
		header := attr.attr(types.AttributeKeyHeaderHolder)
		return ctx.OnDownStreamReceived(header.(Header), body, nil)
	}

	return types.ActionContinue
}

//export proxy_on_request_trailers
func proxyOnRequestTrailers(contextID uint32, numTrailers int) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_trailers")
	}
	this.setActiveContextID(contextID)

	attr := ctx.(attribute)

	var trailer Header
	if numTrailers > 0 {
		trailers, err := getHttpRequestTrailers()
		if err != nil {
			log.Errorf("failed to get request trailer: %v", err)
			return types.ActionContinue
		}
		trailer = &CommonHeader{m: trailers}
		// update context header
		attr.set(types.AttributeKeyTrailerHolder, trailer)
	}

	header := attr.attr(types.AttributeKeyHeaderHolder)
	body := attr.attr(types.AttributeKeyBufferHolder)

	return ctx.OnDownStreamReceived(header.(Header), body.(Buffer), trailer)
}

//export proxy_on_response_headers
func proxyOnResponseHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	var header Header
	if numHeaders > 0 {
		hs, err := getHttpResponseHeaders()
		if err != nil {
			log.Errorf("failed to get response headers: %v", err)
			return types.ActionContinue
		}
		header = &CommonHeader{m: hs}
		attr := ctx.(attribute)
		// update context header
		attr.set(types.AttributeKeyHeaderHolder, header)
	}

	if endOfStream {
		return ctx.OnUpstreamReceived(header, NewBuffer(0), nil)
	}

	return types.ActionContinue
}

//export proxy_on_response_body
func proxyOnResponseBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	attr := ctx.(attribute)

	var body Buffer
	if bodySize > 0 {
		bodyBytes, err := getHttpResponseBody(0, bodySize)
		if err != nil {
			log.Errorf("failed to get response body: %v", err)
			return types.ActionContinue
		}

		body = WrapBuffer(bodyBytes)
		// update context body buffer
		attr.set(types.AttributeKeyBufferHolder, body)
	}

	if endOfStream {
		header := attr.attr(types.AttributeKeyHeaderHolder)
		return ctx.OnUpstreamReceived(header.(Header), body, nil)
	}

	return types.ActionContinue
}

//export proxy_on_response_trailers
func proxyOnResponseTrailers(contextID uint32, numTrailers int) types.Action {
	ctx, ok := this.filterStreams[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	this.setActiveContextID(contextID)

	attr := ctx.(attribute)

	var trailer Header
	if numTrailers > 0 {
		trailers, err := getHttpResponseTrailers()
		if err != nil {
			log.Errorf("failed to get request trailer: %v", err)
			return types.ActionContinue
		}
		trailer = &CommonHeader{m: trailers}
		// update context header
		attr.set(types.AttributeKeyTrailerHolder, trailer)
	}

	header := attr.attr(types.AttributeKeyHeaderHolder)
	body := attr.attr(types.AttributeKeyBufferHolder)

	return ctx.OnUpstreamReceived(header.(Header), body.(Buffer), trailer)
}
