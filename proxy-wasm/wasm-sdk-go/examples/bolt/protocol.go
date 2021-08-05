package bolt

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy"
)

/**
 * Request command protocol for v1
 * 0     1     2           4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |proto| type| cmd code  |ver2 |   requestID           |boltCodec|        Timeout        |  classLen |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 * |headerLen  | contentLen            |                             ... ...                       |
 * +-----------+-----------+-----------+                                                                                               +
 * |               className + header  + content  bytes                                            |
 * +                                                                                               +
 * |                               ... ...                                                         |
 * +-----------------------------------------------------------------------------------------------+
 *
 * proto: code for protocol
 * type: request/response/request oneway
 * cmd code: code for remoting command
 * ver2:version for remoting command
 * requestID: id of request
 * boltCodec: code for boltCodec
 * headerLen: length of header
 * contentLen: length of content
 *
 * Response command protocol for v1
 * 0     1     2     3     4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |proto| type| cmd code  |ver2 |   requestID           |boltCodec|resp Status|  classLen |headerLen  |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 * | contentLen            |                  ... ...                                              |
 * +-----------------------+                                                                       +
 * |                         className + header  + content  bytes                                  |
 * +                                                                                               +
 * |                               ... ...                                                         |
 * +-----------------------------------------------------------------------------------------------+
 * response Status: response Status
 */

type boltProtocol struct {
	boltCodec
	proxy.DefaultOptions
}

func NewBoltProtocol() proxy.Protocol {
	return &boltProtocol{}
}

// types.ProtocolName
func (proto *boltProtocol) Name() string {
	return ProtocolName
}

func (proto *boltProtocol) Codec() proxy.Codec {
	return &proto.boltCodec
}

// heartbeat
func (proto *boltProtocol) KeepAlive(requestId uint64) proxy.Request {
	return &Request{
		RpcHeader: RpcHeader{
			Protocol:  ProtocolCode,
			CmdType:   CmdTypeRequest,
			CmdCode:   CmdCodeHeartbeat,
			Version:   1,
			RequestId: uint32(requestId),
			Codec:     Hessian2Serialize,
		},
		Timeout: -1,
	}
}

func (proto *boltProtocol) ReplyKeepAlive(request proxy.Request) proxy.Response {
	return &Response{
		RpcHeader: RpcHeader{
			Protocol:  ProtocolCode,
			CmdType:   CmdTypeResponse,
			CmdCode:   CmdCodeHeartbeat,
			Version:   ProtocolVersion,
			RequestId: uint32(request.CommandId()),
			Codec:     Hessian2Serialize,
		},
		Status: 0,
	}
}

// Hijacker
func (proto *boltProtocol) Hijack(request proxy.Request, code uint32) proxy.Response {
	return &Response{
		RpcHeader: RpcHeader{
			Protocol:  ProtocolCode,
			CmdType:   CmdTypeResponse,
			CmdCode:   CmdCodeRpcResponse,
			Version:   ProtocolVersion,
			RequestId: uint32(request.CommandId()), // this would be overwrite by stream layer
			Codec:     Hessian2Serialize,           //todo: read default boltCodec from config
		},
		Status: uint16(code),
	}
}
