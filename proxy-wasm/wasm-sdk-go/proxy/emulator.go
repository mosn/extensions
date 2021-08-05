package proxy

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/syscall"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
	stdout "log"
	"sync"
)

type HostEmulator interface {
	// release emulator resource and reset state
	Done()

	// Root
	StartVM()
	StartPlugin()
	FinishVM()

	GetLogs(level types.LogLevel) []string
	GetTickPeriod() uint32
	Tick()

	// protocol L7 level
	NewProtocolContext() (contextID uint32)
	Decode(contextID uint32, data Buffer) (Command, error)
	Encode(contextID uint32, cmd Command) (Buffer, error)
	// heartbeat
	KeepAlive(contextID uint32, requestId uint64) Request
	ReplyKeepAlive(contextID uint32, request Request) Response
	// hijacker
	Hijack(contextID uint32, request Request, code uint32) Response
	// easy for test only
	CurrentStreamId() uint64
	// finish protocol pipeline
	CompleteProtocolContext(contextID uint32)

	// filter L7 level
	NewFilterContext() (contextID uint32)
	PutRequestHeaders(contextID uint32, headers map[string]string, endOfStream bool)
	PutRequestBody(contextID uint32, body []byte, endOfStream bool)
	PutRequestTrailers(contextID uint32, headers map[string]string)

	GetRequestHeaders(contextID uint32) (headers map[string]string)
	GetRequestBody(contextID uint32) []byte

	PutResponseHeaders(contextID uint32, headers map[string]string, endOfStream bool)
	PutResponseBody(contextID uint32, body []byte, endOfStream bool)
	PutResponseTrailers(contextID uint32, headers map[string]string)

	GetResponseHeaders(contextID uint32) (headers map[string]string)
	GetResponseBody(contextID uint32) []byte

	GetCurrentStreamAction(contextID uint32) types.Action
	CompleteFilterContext(contextID uint32)

	// network
}

var (
	hostMux       = sync.Mutex{}
	nextContextID = RootContextID + 1
)

func getNextContextID() (ret uint32) {
	ret = nextContextID
	nextContextID++
	return
}

type hostEmulator struct {
	*rootEmulator
	*filterEmulator
	*protocolEmulator
	*networkEmulator

	effectiveContextID uint32
}

func NewHostEmulator(opt *Option) HostEmulator {
	root := newRootEmulator(opt.pluginConfiguration, opt.vmConfiguration)
	filter := newFilterEmulator()
	protocol := newProtocolEmulator()
	network := newNetworkEmulator()
	emulator := &hostEmulator{
		root,
		filter,
		protocol,
		network,
		0,
	}

	hostMux.Lock() // acquire the lock of host emulation

	syscall.RegisterWasmHost(emulator)

	// set up state
	SetNewRootContext(opt.newRootContext)
	SetNewStreamContext(opt.newStreamContext)
	SetNewFilterContext(opt.newFilterContext)
	SetNewProtocolContext(opt.newProtocolContext)

	// create root context: TODO: support multiple root contexts
	proxyOnContextCreate(RootContextID, 0)

	return emulator
}

// impl HostEmulator
func (*hostEmulator) Done() {
	defer hostMux.Unlock()
	defer VMStateReset()
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	switch bt {
	case types.BufferTypePluginConfiguration, types.BufferTypeVMConfiguration, types.BufferTypeHttpCallResponseBody:
		return h.rootEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	case types.BufferTypeHttpRequestBody, types.BufferTypeHttpResponseBody:
		return h.filterEmulatorProxyGetBufferBytes(bt, start, maxSize, returnBufferData, returnBufferSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxySetBufferBytes(bt types.BufferType, start int, maxSize int, bufferData *byte, bufferSize int) types.Status {
	switch bt {
	case types.BufferTypeHttpRequestBody, types.BufferTypeHttpResponseBody:
		return h.filterEmulatorProxySetBufferBytes(bt, start, maxSize, bufferData, bufferSize)
	case types.BufferTypeDecodeData, types.BufferTypeEncodeData:
		{
			return h.protocolEmulatorProxySetBufferBytes(bt, start, maxSize, bufferData, bufferSize)
		}
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxyGetHeaderMapValue(mapType types.MapType, keyData *byte,
	keySize int, returnValueData **byte, returnValueSize *int) types.Status {
	switch mapType {
	case types.MapTypeHttpRequestHeaders, types.MapTypeHttpResponseHeaders,
		types.MapTypeHttpRequestTrailers, types.MapTypeHttpResponseTrailers:
		return h.filterEmulatorProxyGetHeaderMapValue(mapType, keyData,
			keySize, returnValueData, returnValueSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxyGetHeaderMapPairs(mapType types.MapType, returnValueData **byte,
	returnValueSize *int) types.Status {
	switch mapType {
	case types.MapTypeHttpRequestHeaders, types.MapTypeHttpResponseHeaders,
		types.MapTypeHttpRequestTrailers, types.MapTypeHttpResponseTrailers:
		return h.filterEmulatorProxyGetHeaderMapPairs(mapType, returnValueData, returnValueSize)
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxySetEffectiveContext(contextID uint32) types.Status {
	h.effectiveContextID = contextID
	return types.StatusOK
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxySetProperty(*byte, int, *byte, int) types.Status {
	panic("unimplemented")
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxyGetProperty(*byte, int, **byte, *int) types.Status {
	stdout.Printf("ProxyGetProperty not implemented in the host emulator yet")
	return types.StatusOK
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxyCloseStream(streamType types.StreamType) types.Status {
	stdout.Printf("ProxyCloseStream not implemented in the host emulator yet")
	return types.StatusOK
}

// impl syscall.WasmHost
func (h *hostEmulator) ProxyDone() types.Status {
	stdout.Printf("ProxyDone not implemented in the host emulator yet")
	return types.StatusOK
}
