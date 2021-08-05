package proxy

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/syscall"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
)

func getPluginConfiguration(size int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypePluginConfiguration, 0, size)
	return buf, types.StatusToError(status)
}

func getVMConfiguration(size int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypeVMConfiguration, 0, size)
	return buf, types.StatusToError(status)
}

func SetTickPeriodMilliSeconds(millSec uint32) error {
	return types.StatusToError(syscall.ProxySetTickPeriodMilliseconds(millSec))
}

func getDownStreamData(start, maxSize int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypeDownstreamData, start, maxSize)
	return buf, types.StatusToError(status)
}

func getUpstreamData(start, maxSize int) ([]byte, error) {
	buf, status := getBuffer(types.BufferTypeUpstreamData, start, maxSize)
	return buf, types.StatusToError(status)
}

func getHttpRequestHeaders() (map[string]string, error) {
	headers, status := getMap(types.MapTypeHttpRequestHeaders)
	return headers, types.StatusToError(status)
}

func setHttpRequestHeaders(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestHeaders, headers))
}

func getHttpRequestHeader(key string) (string, error) {
	header, status := getMapValue(types.MapTypeHttpRequestHeaders, key)
	return header, types.StatusToError(status)
}

func removeHttpRequestHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestHeaders, key))
}

func setHttpRequestHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

func addHttpRequestHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

func getHttpRequestBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpRequestBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func setHttpRequestBody(body []byte) error {
	var buff *byte
	if len(body) != 0 {
		buff = &body[0]
	}
	status := syscall.ProxySetBufferBytes(types.BufferTypeHttpRequestBody, 0, len(body), buff, len(body))
	return types.StatusToError(status)
}

func setDecodeBuffer(body []byte) error {
	var buff *byte
	if len(body) != 0 {
		buff = &body[0]
	}
	status := syscall.ProxySetBufferBytes(types.BufferTypeDecodeData, 0, len(body), buff, len(body))
	return types.StatusToError(status)
}

func setEncodeBuffer(body []byte) error {
	var buff *byte
	if len(body) != 0 {
		buff = &body[0]
	}
	status := syscall.ProxySetBufferBytes(types.BufferTypeEncodeData, 0, len(body), buff, len(body))
	return types.StatusToError(status)
}

func getHttpRequestTrailers() (map[string]string, error) {
	trailers, status := getMap(types.MapTypeHttpRequestTrailers)
	return trailers, types.StatusToError(status)
}

func setHttpRequestTrailers(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestTrailers, headers))
}

func getHttpRequestTrailer(key string) (string, error) {
	trailer, status := getMapValue(types.MapTypeHttpRequestTrailers, key)
	return trailer, types.StatusToError(status)
}

func removeHttpRequestTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestTrailers, key))
}

func setHttpRequestTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

func addHttpRequestTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

func resumeHttpRequest() error {
	return types.StatusToError(syscall.ProxyContinueStream(types.StreamTypeRequest))
}

func getHttpResponseHeaders() (map[string]string, error) {
	headers, status := getMap(types.MapTypeHttpResponseHeaders)
	return headers, types.StatusToError(status)
}

func setHttpResponseHeaders(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseHeaders, headers))
}

func getHttpResponseHeader(key string) (string, error) {
	header, status := getMapValue(types.MapTypeHttpResponseHeaders, key)
	return header, types.StatusToError(status)
}

func removeHttpResponseHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseHeaders, key))
}

func setHttpResponseHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

func addHttpResponseHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

func getHttpResponseBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpResponseBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func setHttpResponseBody(body []byte) error {
	var buf *byte
	if len(body) != 0 {
		buf = &body[0]
	}
	st := syscall.ProxySetBufferBytes(types.BufferTypeHttpResponseBody, 0, len(body), buf, len(body))
	return types.StatusToError(st)
}

func getHttpResponseTrailers() (map[string]string, error) {
	trailers, status := getMap(types.MapTypeHttpResponseTrailers)
	return trailers, types.StatusToError(status)
}

func setHttpResponseTrailers(headers map[string]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseTrailers, headers))
}

func getHttpResponseTrailer(key string) (string, error) {
	trailer, status := getMapValue(types.MapTypeHttpResponseTrailers, key)
	return trailer, types.StatusToError(status)
}

func removeHttpResponseTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseTrailers, key))
}

func setHttpResponseTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

func addHttpResponseTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

func resumeHttpResponse() error {
	return types.StatusToError(syscall.ProxyContinueStream(types.StreamTypeResponse))
}

func getProperty(path []string) ([]byte, error) {
	var ret *byte
	var retSize int
	raw := EncodePropertyPath(path)

	err := types.StatusToError(syscall.ProxyGetProperty(&raw[0], len(raw), &ret, &retSize))
	if err != nil {
		return nil, err
	}

	return parseByteSlice(ret, retSize), nil
}

func setProperty(path string, data []byte) error {
	return types.StatusToError(syscall.ProxySetProperty(
		stringBytePtr(path), len(path), &data[0], len(data),
	))
}

func setMap(mapType types.MapType, headers map[string]string) types.Status {
	encodedBytes := EncodeMap(headers)
	hp := &encodedBytes[0]
	hl := len(encodedBytes)
	return syscall.ProxySetHeaderMapPairs(mapType, hp, hl)
}

func getMapValue(mapType types.MapType, key string) (string, types.Status) {
	var rvs int
	var raw *byte
	if st := syscall.ProxyGetHeaderMapValue(mapType, stringBytePtr(key), len(key), &raw, &rvs); st != types.StatusOK {
		return "", st
	}

	ret := parseString(raw, rvs)
	return ret, types.StatusOK
}

func removeMapValue(mapType types.MapType, key string) types.Status {
	return syscall.ProxyRemoveHeaderMapValue(mapType, stringBytePtr(key), len(key))
}

func setMapValue(mapType types.MapType, key, value string) types.Status {
	return syscall.ProxyReplaceHeaderMapValue(mapType, stringBytePtr(key), len(key), stringBytePtr(value), len(value))
}

func addMapValue(mapType types.MapType, key, value string) types.Status {
	return syscall.ProxyAddHeaderMapValue(mapType, stringBytePtr(key), len(key), stringBytePtr(value), len(value))
}

func getMap(mapType types.MapType) (map[string]string, types.Status) {
	var rvs int
	var raw *byte

	status := syscall.ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if status != types.StatusOK {
		return nil, status
	}

	bs := parseByteSlice(raw, rvs)
	return DecodeMap(bs), types.StatusOK
}

func getBuffer(bufType types.BufferType, start, maxSize int) ([]byte, types.Status) {
	var buffer *byte
	var len int
	switch status := syscall.ProxyGetBufferBytes(bufType, start, maxSize, &buffer, &len); status {
	case types.StatusOK:
		return parseByteSlice(buffer, len), status
	default:
		return nil, status
	}
}
