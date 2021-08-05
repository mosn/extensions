package syscall

import "github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"

type WasmHost interface {
	ProxyLog(logLevel types.LogLevel, buffer *byte, len int) types.Status
	ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) types.Status
	ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) types.Status
	ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status
	ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status
	ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status
	ProxyContinueStream(streamType types.StreamType) types.Status
	ProxyCloseStream(streamType types.StreamType) types.Status
	ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status
	ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status
	ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status
	ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status
	ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status
	ProxySetTickPeriodMilliseconds(period uint32) types.Status
	ProxySetEffectiveContext(contextID uint32) types.Status
	ProxyDone() types.Status
}

var proxyHost WasmHost

func RegisterWasmHost(wasmHost WasmHost) {
	proxyHost = wasmHost
}

type DefaultWasmHost struct {
}

var defaultHost WasmHost = &DefaultWasmHost{}

func (h *DefaultWasmHost) ProxyLog(logLevel types.LogLevel, buffer *byte, len int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyAddHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyReplaceHeaderMapValue(mapType types.MapType, keyData *byte, keySize int, valueData *byte, valueSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyContinueStream(streamType types.StreamType) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyCloseStream(streamType types.StreamType) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyRemoveHeaderMapValue(mapType types.MapType, keyData *byte, keySize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte, returnValueSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetHeaderMapPairs(mapType types.MapType, mapData *byte, mapSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxySetEffectiveContext(contextID uint32) types.Status {
	return types.StatusOK
}

func (h *DefaultWasmHost) ProxyDone() types.Status { return types.StatusOK }

func wasmHost() WasmHost {
	if proxyHost == nil {
		return defaultHost
	}
	return proxyHost
}
