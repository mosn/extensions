package proxy

import (
	"context"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
	"sync/atomic"
)

type Protocol interface {
	Name() string
	Codec() Codec
	KeepAlive
	Hijacker
	Options
}

type Codec interface {
	Decode(ctx context.Context, data Buffer) (Command, error)
	Encode(ctx context.Context, cmd Command) (Buffer, error)
}

// Command base request or response command
type Command interface {
	// Header get the data exchange header, maybe return nil.
	GetHeader() Header
	// GetData return the full message buffer, the protocol header is not included
	GetData() Buffer
	// SetData update the full message buffer, the protocol header is not included
	SetData(data Buffer)
	// IsHeartbeat check if the request is a heartbeat request
	IsHeartbeat() bool
	// CommandId get command id
	CommandId() uint64
	// SetCommandId update command id
	// In upstream, because of connection multiplexing,
	// the id of downstream needs to be replaced with id of upstream
	// blog: https://mosn.io/blog/posts/multi-protocol-deep-dive/#%E5%8D%8F%E8%AE%AE%E6%89%A9%E5%B1%95%E6%A1%86%E6%9E%B6
	SetCommandId(id uint64)
}

type Request interface {
	Command
	// IsOneWay Check that the request does not care about the response
	IsOneWay() bool
	GetTimeout() uint32 // request timeout
}

type Response interface {
	Command
	GetStatus() uint32 // response status
}

type KeepAlive interface {
	KeepAlive(requestId uint64) Request
	ReplyKeepAlive(request Request) Response
}

type Hijacker interface {
	// Hijack allows sidecar to hijack requests
	Hijack(request Request, code uint32) Response
}

// Options Features required for a particular host
type Options interface {
	// default Multiplex
	PoolMode() types.PoolMode
	// EnableWorkerPool same meaning as EnableWorkerPool in types.StreamConnection
	EnableWorkerPool() bool
	// EnableGenerateRequestID check if the protocol requires custom streamId.
	// If need to generate, you should override the GenerateRequestID method implementation
	EnableGenerateRequestID() bool
	// GenerateRequestID generate a request id for stream to combine stream request && response
	// use connection param as base
	GenerateRequestID(*uint64) uint64
}

var options = &DefaultOptions{}

type DefaultOptions struct {
}

func (o *DefaultOptions) PoolMode() types.PoolMode {
	return types.Multiplex
}

func (o *DefaultOptions) EnableWorkerPool() bool {
	return true
}

func (o *DefaultOptions) EnableGenerateRequestID() bool {
	return false
}

func (o *DefaultOptions) GenerateRequestID(v *uint64) uint64 {
	return atomic.AddUint64(v, 1)
}

const (
	ResponseStatusSuccess                    uint16 = 0  // 0x00 response status
	ResponseStatusError                      uint16 = 1  // 0x01
	ResponseStatusServerException            uint16 = 2  // 0x02
	ResponseStatusUnknown                    uint16 = 3  // 0x03
	ResponseStatusServerThreadPoolBusy       uint16 = 4  // 0x04
	ResponseStatusErrorComm                  uint16 = 5  // 0x05
	ResponseStatusNoProcessor                uint16 = 6  // 0x06
	ResponseStatusTimeout                    uint16 = 7  // 0x07
	ResponseStatusClientSendError            uint16 = 8  // 0x08
	ResponseStatusCodecException             uint16 = 9  // 0x09
	ResponseStatusConnectionClosed           uint16 = 16 // 0x10
	ResponseStatusServerSerialException      uint16 = 17 // 0x11
	ResponseStatusServerDeserializeException uint16 = 18 // 0x12

	CodecExceptionCode       = 0
	UnknownCode              = 2
	DeserializeExceptionCode = 3
	SuccessCode              = 200
	PermissionDeniedCode     = 403
	RouterUnavailableCode    = 404
	InternalErrorCode        = 500
	NoHealthUpstreamCode     = 502
	UpstreamOverFlowCode     = 503
	TimeoutExceptionCode     = 504
	LimitExceededCode        = 509
)

func Mapping(httpStatusCode int32) uint32 {
	switch httpStatusCode {
	case SuccessCode:
		return uint32(ResponseStatusSuccess)
	case RouterUnavailableCode:
		return uint32(ResponseStatusNoProcessor)
	case NoHealthUpstreamCode:
		return uint32(ResponseStatusConnectionClosed)
	case UpstreamOverFlowCode:
		return uint32(ResponseStatusServerThreadPoolBusy)
	case CodecExceptionCode:
		//Decode or Encode Error
		return uint32(ResponseStatusCodecException)
	case DeserializeExceptionCode:
		//Hessian Exception
		return uint32(ResponseStatusServerDeserializeException)
	case TimeoutExceptionCode:
		//Response Timeout
		return uint32(ResponseStatusTimeout)
	default:
		return uint32(ResponseStatusUnknown)
	}
}
