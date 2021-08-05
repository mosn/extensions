package proxy

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
	stdout "log"
)

type (
	filterEmulator struct {
		filterStreams map[uint32]*filterStreamState
	}
	filterStreamState struct {
		requestHeaders    map[string]string
		requestBody       []byte
		requestTrailers   map[string]string
		responseHeaders   map[string]string
		responseBody      []byte
		responseTrailers  map[string]string
		action            types.Action
		sentLocalResponse *LocalHttpResponse
	}
	LocalHttpResponse struct {
		StatusCode       uint32
		StatusCodeDetail string
		Data             []byte
		Headers          map[string]string
		GRPCStatus       int32
	}
)

func newFilterEmulator() *filterEmulator {
	host := &filterEmulator{filterStreams: map[uint32]*filterStreamState{}}
	return host
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) filterEmulatorProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]
	var buf []byte
	switch bt {
	case types.BufferTypeHttpRequestBody:
		buf = stream.requestBody
	case types.BufferTypeHttpResponseBody:
		buf = stream.responseBody
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if start >= len(buf) {
		stdout.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return types.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return types.StatusOK
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) filterEmulatorProxySetBufferBytes(bt types.BufferType, start int, maxSize int,
	bufferData *byte, bufferSize int) types.Status {
	body := parseByteSlice(bufferData, bufferSize)
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]
	switch bt {
	case types.BufferTypeHttpRequestBody:
		stream.requestBody = body
	case types.BufferTypeHttpResponseBody:
		stream.responseBody = body
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
	return types.StatusOK
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) filterEmulatorProxyGetHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	key := parseString(keyData, keySize)
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]

	var headers map[string]string
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		headers = stream.requestHeaders
	case types.MapTypeHttpResponseHeaders:
		headers = stream.responseHeaders
	case types.MapTypeHttpRequestTrailers:
		headers = stream.requestTrailers
	case types.MapTypeHttpResponseTrailers:
		headers = stream.responseTrailers
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	for k, v := range headers {
		if k == key {
			value := []byte(v)
			*returnValueData = &value[0]
			*returnValueSize = len(value)
			return types.StatusOK
		}
	}

	return types.StatusNotFound
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) types.Status {

	key := parseString(keyData, keySize)
	value := parseString(valueData, valueSize)
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = updateMapValue(stream.requestHeaders, key, value)
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = updateMapValue(stream.responseHeaders, key, value)
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = updateMapValue(stream.requestTrailers, key, value)
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = updateMapValue(stream.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}

	return types.StatusOK
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, valueData *byte, valueSize int) types.Status {
	key := parseString(keyData, keySize)
	value := parseString(valueData, valueSize)
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = updateMapValue(stream.requestHeaders, key, value)
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = updateMapValue(stream.responseHeaders, key, value)
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = updateMapValue(stream.requestTrailers, key, value)
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = updateMapValue(stream.responseTrailers, key, value)
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	key := parseString(keyData, keySize)
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = removeMapValue0(stream.requestHeaders, key)
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = removeMapValue0(stream.responseHeaders, key)
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = removeMapValue0(stream.requestTrailers, key)
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = removeMapValue0(stream.responseTrailers, key)
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

// impl syscall.WasmHost: delegated from hostEmulator
func (h *filterEmulator) filterEmulatorProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte,
	returnValueSize *int) types.Status {
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]

	var m []byte
	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		m = EncodeMap(stream.requestHeaders)
	case types.MapTypeHttpResponseHeaders:
		m = EncodeMap(stream.responseHeaders)
	case types.MapTypeHttpRequestTrailers:
		m = EncodeMap(stream.requestTrailers)
	case types.MapTypeHttpResponseTrailers:
		m = EncodeMap(stream.responseTrailers)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	*returnValueData = &m[0]
	*returnValueSize = len(m)
	return types.StatusOK
}

func (h *filterEmulator) ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	m := DecodeMap(parseByteSlice(mapData, mapSize))
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]

	switch mapType {
	case types.MapTypeHttpRequestHeaders:
		stream.requestHeaders = m
	case types.MapTypeHttpResponseHeaders:
		stream.responseHeaders = m
	case types.MapTypeHttpRequestTrailers:
		stream.requestTrailers = m
	case types.MapTypeHttpResponseTrailers:
		stream.responseTrailers = m
	default:
		panic("unimplemented")
	}
	return types.StatusOK
}

func (h *filterEmulator) ProxyContinueStream(types.StreamType) types.Status {
	active := VMStateGetActiveContextID()
	stream := h.filterStreams[active]
	stream.action = types.ActionContinue
	return types.StatusOK
}

// impl HostEmulator
func (h *filterEmulator) NewFilterContext() (contextID uint32) {
	contextID = getNextContextID()
	proxyOnContextCreate(contextID, RootContextID)
	h.filterStreams[contextID] = &filterStreamState{action: types.ActionContinue}
	return
}

// impl HostEmulator
func (h *filterEmulator) PutRequestHeaders(contextID uint32, headers map[string]string, endOfStream bool) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestHeaders = headers
	cs.action = proxyOnRequestHeaders(contextID, len(headers), endOfStream)
}

// impl HostEmulator
func (h *filterEmulator) PutRequestBody(contextID uint32, body []byte, endOfStream bool) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestBody = body
	cs.action = proxyOnRequestBody(contextID, len(body), endOfStream)
}

// impl HostEmulator
func (h *filterEmulator) PutRequestTrailers(contextID uint32, headers map[string]string) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.requestTrailers = headers
	cs.action = proxyOnRequestTrailers(contextID, len(headers))
}

// impl HostEmulator
func (h *filterEmulator) GetRequestHeaders(contextID uint32) (headers map[string]string) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}
	return cs.requestHeaders
}

// impl HostEmulator
func (h *filterEmulator) GetRequestBody(contextID uint32) []byte {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}
	return cs.requestBody
}

// impl HostEmulator
func (h *filterEmulator) PutResponseHeaders(contextID uint32, headers map[string]string, endOfStream bool) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseHeaders = headers
	cs.action = proxyOnResponseHeaders(contextID, len(headers), endOfStream)
}

// impl HostEmulator
func (h *filterEmulator) PutResponseBody(contextID uint32, body []byte, endOfStream bool) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseBody = body
	cs.action = proxyOnResponseBody(contextID, len(body), endOfStream)
}

// impl HostEmulator
func (h *filterEmulator) PutResponseTrailers(contextID uint32, headers map[string]string) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}

	cs.responseTrailers = headers
	cs.action = proxyOnRequestTrailers(contextID, len(headers))
}

// impl HostEmulator
func (h *filterEmulator) GetResponseHeaders(contextID uint32) (headers map[string]string) {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}
	return cs.responseHeaders
}

// impl HostEmulator
func (h *filterEmulator) GetResponseBody(contextID uint32) []byte {
	cs, ok := h.filterStreams[contextID]
	if !ok {
		stdout.Fatalf("invalid context id: %d", contextID)
	}
	return cs.responseBody
}

// impl HostEmulator
func (h *filterEmulator) CompleteFilterContext(contextID uint32) {
	proxyOnLog(contextID)
	proxyOnDone(contextID)
	proxyOnDelete(contextID)
}

// impl HostEmulator
func (h *filterEmulator) GetCurrentStreamAction(contextID uint32) types.Action {
	stream, ok := h.filterStreams[contextID]
	if !ok {
		log.Fatalf("invalid context id: %d", contextID)
	}
	return stream.action
}

func updateMapValue(base map[string]string, key, value string) map[string]string {
	base[key] = value
	return base
}

func removeMapValue0(base map[string]string, key string) map[string]string {
	delete(base, key)
	return base
}
