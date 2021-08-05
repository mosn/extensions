package types

import (
	"errors"
	"strconv"
)

// Action tell the host what action should be triggered
type Action uint32

const (
	ActionContinue Action = 0
	ActionPause    Action = 1
)

// Status
type Status uint32

const (
	StatusOK              Status = 0
	StatusNotFound        Status = 1
	StatusBadArgument     Status = 2
	StatusEmpty           Status = 7
	StatusCasMismatch     Status = 8
	StatusInternalFailure Status = 10
	StatusNeedMoreData    Status = 99
)

type StreamType uint32

const (
	StreamTypeRequest  StreamType = 0
	StreamTypeResponse StreamType = 1
)

type ContextType uint32

const (
	VmContext     ContextType = 1
	PluginContext ContextType = 2
	StreamContext ContextType = 3
	HttpContext   ContextType = 4
)

type BufferType uint32

const (
	BufferTypeHttpRequestBody      BufferType = 0
	BufferTypeHttpResponseBody     BufferType = 1
	BufferTypeDownstreamData       BufferType = 2
	BufferTypeUpstreamData         BufferType = 3
	BufferTypeHttpCallResponseBody BufferType = 4
	BufferTypeGrpcReceiveBuffer    BufferType = 5
	BufferTypeVMConfiguration      BufferType = 6
	BufferTypePluginConfiguration  BufferType = 7
	BufferTypeCallData             BufferType = 8
	BufferTypeDecodeData           BufferType = 13
	BufferTypeEncodeData           BufferType = 14
)

type MapType uint32

const (
	MapTypeHttpRequestHeaders       MapType = 0
	MapTypeHttpRequestTrailers      MapType = 1
	MapTypeHttpResponseHeaders      MapType = 2
	MapTypeHttpResponseTrailers     MapType = 3
	MapTypeHttpCallResponseHeaders  MapType = 6
	MapTypeHttpCallResponseTrailers MapType = 7
	MapTypeRpcRequestHeaders        MapType = 31
	MapTypeRpcRequestTrailers       MapType = 32
	MapTypeRpcResponseHeaders       MapType = 33
	MapTypeRpcResponseTrailers      MapType = 34
)

// PeerType
type PeerType uint32

const (
	Local  PeerType = 1
	Remote PeerType = 2
)

// LogLevel proxy log level
type LogLevel uint32

const (
	LogLevelTrace    LogLevel = 0
	LogLevelDebug    LogLevel = 1
	LogLevelInfo     LogLevel = 2
	LogLevelWarn     LogLevel = 3
	LogLevelError    LogLevel = 4
	LogLevelCritical LogLevel = 5
	LogLevelFatal    LogLevel = 6
	LogLevelMax      LogLevel = 7
)

const (
	traceText    = "trace"
	debugText    = "debug"
	infoText     = "info"
	warnText     = "warn"
	errorText    = "error"
	fatalText    = "fatal"
	criticalText = "critical"
)

func (level LogLevel) String() string {
	switch level {
	case LogLevelTrace:
		return traceText
	case LogLevelDebug:
		return debugText
	case LogLevelInfo:
		return infoText
	case LogLevelWarn:
		return warnText
	case LogLevelError:
		return errorText
	case LogLevelCritical:
		return criticalText
	case LogLevelFatal:
		return fatalText
	default:
		panic("unsupported log level")
	}
}

type ExtensionType int

const (
	VmContextFilter     ExtensionType = 1
	PluginContextFilter ExtensionType = 2
	StreamContextFilter ExtensionType = 3
	HttpContextFilter   ExtensionType = 4
)

var (
	ErrorStatusNotFound    = errors.New("error status returned by host: not found")
	ErrorStatusBadArgument = errors.New("error status returned by host: bad argument")
	ErrorStatusEmpty       = errors.New("error status returned by host: empty")
	ErrorStatusCasMismatch = errors.New("error status returned by host: cas mismatch")
	ErrorInternalFailure   = errors.New("error status returned by host: internal failure")
)

//go:inline
func StatusToError(status Status) error {
	switch status {
	case StatusOK:
		return nil
	case StatusNotFound:
		return ErrorStatusNotFound
	case StatusBadArgument:
		return ErrorStatusBadArgument
	case StatusEmpty:
		return ErrorStatusEmpty
	case StatusCasMismatch:
		return ErrorStatusCasMismatch
	case StatusInternalFailure:
		return ErrorInternalFailure
	}
	return errors.New("unknown status code: " + strconv.Itoa(int(status)))
}

// context impl
type AttributeKey string

const (
	AttributeKeyStreamID      AttributeKey = "stream_id"
	AttributeKeyListenerType               = "listener_type"
	AttributeKeyHeaderHolder               = "header"
	AttributeKeyBufferHolder               = "buffer"
	AttributeKeyTrailerHolder              = "trailer"
	AttributeKeyDecodeCommand              = "decode_command"
	AttributeKeyEncodeCommand              = "encode_command"
	AttributeKeyEncodedBuffer              = "encoded_buffer"
)

// PoolMode is whether PingPong or multiplex
type PoolMode int

const (
	PingPong PoolMode = iota
	Multiplex
	TCP
)

const (
	ResponseType       = 0
	RequestType        = 1
	RequestOneWayType  = 2
	UnKnownRpcFlagType = "unknown protocol flag type"
)

// ContextKey type
type ContextKey int

// Context key types(built-in)
const (
	ContextKeyStreamID ContextKey = iota
	ContextKeyConnectionID
	ContextKeyListenerPort
	ContextKeyListenerName
	ContextKeyListenerType
	ContextKeyListenerStatsNameSpace
	ContextKeyNetworkFilterChainFactories
	ContextKeyStreamFilterChainFactories
	ContextKeyBufferPoolCtx
	ContextKeyAccessLogs
	ContextOriRemoteAddr
	ContextKeyAcceptChan
	ContextKeyAcceptBuffer
	ContextKeyConnectionFd
	ContextSubProtocol
	ContextKeyTraceSpanKey
	ContextKeyActiveSpan
	ContextKeyTraceId
	ContextKeyVariables
	ContextKeyProxyGeneralConfig
	ContextKeyDownStreamProtocol
	ContextKeyConfigDownStreamProtocol
	ContextKeyConfigUpStreamProtocol
	ContextKeyDownStreamHeaders
	ContextKeyDownStreamRespHeaders
	ContextKeyStreamFilterPhase
	ContextKeyEnd
)
