package proxy

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
	stdout "log"
)

type (
	rootEmulator struct {
		logs                        [types.LogLevelMax][]string
		tickPeriod                  uint32
		httpContextIDToCalloutInfos map[uint32][]HttpCalloutAttribute // key: contextID
		httpCalloutIDToContextID    map[uint32]uint32                 // key: calloutID -> contextID
		httpCalloutResponse         map[uint32]HttpCalloutResponse    // key: calloutID
		pluginConfiguration         ConfigMap
		vmConfiguration             ConfigMap
		pluginConfigurationBytes    []byte
		vmConfigurationBytes        []byte
		activeCalloutID             uint32
	}

	HttpCalloutAttribute struct {
		CalloutID uint32
		Upstream  string
		Headers   map[string]string
		Body      []byte
		Trailers  map[string]string
	}

	HttpCalloutResponse struct {
		headers  map[string]string
		body     []byte
		trailers map[string]string
	}
)

func newRootEmulator(pluginConfiguration, vmConfiguration ConfigMap) *rootEmulator {
	host := &rootEmulator{
		httpContextIDToCalloutInfos: map[uint32][]HttpCalloutAttribute{},
		httpCalloutIDToContextID:    map[uint32]uint32{},
		httpCalloutResponse:         map[uint32]HttpCalloutResponse{},
		pluginConfiguration:         pluginConfiguration,
		vmConfiguration:             vmConfiguration,
	}
	if pluginConfiguration != nil {
		host.pluginConfigurationBytes = EncodeMap(pluginConfiguration.ToMap())
	}
	if vmConfiguration != nil {
		host.vmConfigurationBytes = EncodeMap(vmConfiguration.ToMap())
	}

	return host
}

// impl syscall.WasmHost #ProxyLog
func (r *rootEmulator) ProxyLog(logLevel types.LogLevel, messageData *byte, messageSize int) types.Status {
	str := parseString(messageData, messageSize)

	stdout.Printf("proxy_%s_log: %s", logLevel, str)
	r.logs[logLevel] = append(r.logs[logLevel], str)
	return types.StatusOK
}

// impl syscall.WasmHost #ProxySetTickPeriodMilliseconds
func (r *rootEmulator) ProxySetTickPeriodMilliseconds(period uint32) types.Status {
	r.tickPeriod = period
	return types.StatusOK
}

// // impl syscall.WasmHost: delegated from HostEmulator
func (r *rootEmulator) rootEmulatorProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	var buf []byte
	switch bt {
	case types.BufferTypePluginConfiguration:
		buf = r.pluginConfigurationBytes
	case types.BufferTypeVMConfiguration:
		buf = r.vmConfigurationBytes
	case types.BufferTypeHttpCallResponseBody:
		activeID := VMStateGetActiveContextID()
		res, ok := r.httpCalloutResponse[r.activeCalloutID]
		if !ok {
			log.Fatalf("callout response unregistered for %d", activeID)
		}
		buf = res.body
	default:
		panic("unreachable: maybe a bug in this host emulation or SDK")
	}

	if len(buf) == 0 {
		// not config found
		return types.StatusOK
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

// impl HostEmulator
func (r *rootEmulator) GetLogs(level types.LogLevel) []string {
	if level >= types.LogLevelMax {
		log.Fatalf("invalid log level: %d", level)
	}
	return r.logs[level]
}

// impl HostEmulator
func (r *rootEmulator) GetTickPeriod() uint32 {
	return r.tickPeriod
}

const (
	RootContextID uint32 = 1 // TODO: support multiple rootContext
)

// impl HostEmulator
func (r *rootEmulator) Tick() {
	proxyOnTick(RootContextID)
}

// impl HostEmulator
func (r *rootEmulator) GetCalloutAttributesFromContext(contextID uint32) []HttpCalloutAttribute {
	infos := r.httpContextIDToCalloutInfos[contextID]
	return infos
}

// impl HostEmulator
func (r *rootEmulator) StartVM() {
	proxyOnVMStart(RootContextID, len(r.vmConfigurationBytes))
}

// impl HostEmulator
func (r *rootEmulator) StartPlugin() {
	proxyOnPluginStart(RootContextID, len(r.pluginConfigurationBytes))
}

// impl HostEmulator
func (r *rootEmulator) FinishVM() {
	proxyOnDone(RootContextID)
}
