package bolt

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy"
)

const (
	ProtocolName    string = "x-bolt" // protocol
	ProtocolCode    byte   = 1
	ProtocolVersion byte   = 1

	CmdTypeResponse      byte = 0 // cmd type
	CmdTypeRequest       byte = 1
	CmdTypeRequestOneway byte = 2

	CmdCodeHeartbeat   uint16 = 0 // cmd code
	CmdCodeRpcRequest  uint16 = 1
	CmdCodeRpcResponse uint16 = 2

	Hessian2Serialize byte = 1 // serialize

	RequestHeaderLen  int = 22 // protocol header fields length
	ResponseHeaderLen int = 20
	LessLen           int = ResponseHeaderLen // minimal length for decoding

	RequestIdIndex         = 5
	RequestHeaderLenIndex  = 16
	ResponseHeaderLenIndex = 14

	UnKnownCmdType string = "unknown cmd type"
)

type RpcHeader struct {
	Protocol  byte // meta fields
	CmdType   byte
	CmdCode   uint16
	Version   byte
	RequestId uint32
	Codec     byte

	ClassLen   uint16
	HeaderLen  uint16
	ContentLen uint32

	Class string // payload fields

	proxy.CommonHeader
}

// Request is the cmd struct of bolt v1 request
type Request struct {
	RpcHeader

	Timeout    int32
	rawData    []byte // raw data
	rawMeta    []byte // sub slice of raw data, start from protocol code, ends to content length
	rawClass   []byte // sub slice of raw data, class bytes
	rawHeader  []byte // sub slice of raw data, header bytes
	rawContent []byte // sub slice of raw data, content bytes

	Data    proxy.Buffer // wrapper of raw data
	Content proxy.Buffer // wrapper of raw content

	ContentChanged bool // indicate that content changed
}

// Header get the data exchange header, maybe return nil.
func (r *Request) GetHeader() proxy.Header {
	return r
}

// GetData return the complete message byte buffer, including the protocol header
func (r *Request) GetData() proxy.Buffer {
	return r.Content
}

// SetData update the complete message byte buffer, including the protocol header
func (r *Request) SetData(data proxy.Buffer) {
	// judge if the address unchanged, assume that proxy logic will not operate the original Content buffer.
	if r.Content != data {
		r.ContentChanged = true
		r.Content = data
	}
}

// IsHeartbeat check if the request is a heartbeat request
func (r *Request) IsHeartbeat() bool {
	return r.CmdCode == CmdCodeHeartbeat
}

// CommandId get command id
func (r *Request) CommandId() uint64 {
	return uint64(r.RequestId)
}

func (r *Request) SetCommandId(id uint64) {
	r.RequestId = uint32(id)
}

// IsOneWay Check that the request does not care about the response
func (r *Request) IsOneWay() bool {
	return r.CmdType == CmdTypeRequestOneway
}

func (r *Request) GetTimeout() uint32 {
	return uint32(r.Timeout)
}

// Response is the cmd struct of bolt v1 response
type Response struct {
	RpcHeader

	Status uint16

	rawData    []byte // raw data
	rawMeta    []byte // sub slice of raw data, start from protocol code, ends to content length
	rawClass   []byte // sub slice of raw data, class bytes
	rawHeader  []byte // sub slice of raw data, header bytes
	rawContent []byte // sub slice of raw data, content bytes

	Data    proxy.Buffer // wrapper of raw data
	Content proxy.Buffer // wrapper of raw content

	ContentChanged bool // indicate that content changed
}

// Header get the data exchange header, maybe return nil.
func (r *Response) GetHeader() proxy.Header {
	return r
}

// GetData return the complete message byte buffer, including the protocol header
func (r *Response) GetData() proxy.Buffer {
	return r.Content
}

// SetData update the complete message byte buffer, including the protocol header
func (r *Response) SetData(data proxy.Buffer) {
	// judge if the address unchanged, assume that proxy logic will not operate the original Content buffer.
	if r.Content != data {
		r.ContentChanged = true
		r.Content = data
	}
}

// IsHeartbeat check if the request is a heartbeat request
func (r *Response) IsHeartbeat() bool {
	return r.CmdCode == CmdCodeHeartbeat
}

// CommandId get command id
func (r *Response) CommandId() uint64 {
	return uint64(r.RequestId)
}

func (r *Response) SetCommandId(id uint64) {
	r.RequestId = uint32(id)
}

// response Status
func (r *Response) GetStatus() uint32 {
	return uint32(r.Status)
}
