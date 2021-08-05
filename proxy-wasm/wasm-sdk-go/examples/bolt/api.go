package bolt

import "github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy"

// NewRpcRequest is a utility function which build rpc Request object of bolt protocol.
func NewRpcRequest(requestId uint32, headers proxy.Header, data proxy.Buffer) *Request {
	request := &Request{
		RpcHeader: RpcHeader{
			Protocol:  ProtocolCode,
			CmdType:   CmdTypeRequest,
			CmdCode:   CmdCodeRpcRequest,
			Version:   ProtocolVersion,
			RequestId: requestId,
			Codec:     Hessian2Serialize,
		},
		Timeout: -1,
	}

	// set headers
	if headers != nil {
		headers.Range(func(key, value string) bool {
			request.Set(key, value)
			return true
		})
	}

	// set content
	if data != nil {
		request.Content = data
	}
	return request
}

// NewRpcResponse is a utility function which build rpc Response object of bolt protocol.
func NewRpcResponse(requestId uint32, statusCode uint16, headers proxy.Header, data proxy.Buffer) *Response {
	response := &Response{
		RpcHeader: RpcHeader{
			Protocol:  ProtocolCode,
			CmdType:   CmdTypeResponse,
			CmdCode:   CmdCodeRpcResponse,
			Version:   ProtocolVersion,
			RequestId: requestId,
			Codec:     Hessian2Serialize,
		},
		Status: statusCode,
	}

	// set headers
	if headers != nil {
		headers.Range(func(key, value string) bool {
			response.Set(key, value)
			return true
		})
	}

	// set content
	if data != nil {
		response.Content = data
	}
	return response
}
