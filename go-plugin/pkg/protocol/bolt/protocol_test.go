package bolt

import (
	"context"
	"mosn.io/pkg/header"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

func TestProto(t *testing.T) {
	var (
		bp      = BoltProtocol{}
		ctx     = context.TODO()
		payload = "hello world"
		header  = header.CommonHeader{"k": "v"}
		req     = NewRpcRequest(111, header, buffer.NewIoBufferString(payload))
		resp    = NewRpcResponse(111, 0, header, buffer.NewIoBufferString(payload))
	)

	assert.Equal(t, bp.Name(), ProtocolName)

	/////// request
	buf, err := bp.Encode(ctx, req)
	assert.Nil(t, err)

	cmdInter, err := bp.Decode(ctx, buf)
	assert.Nil(t, err)
	cmd, ok := cmdInter.(*Request)
	assert.True(t, ok)
	assert.Equal(t, cmd.Content.String(), payload)
	/////// request end

	/////// heartbeat
	frame := bp.Trigger(context.TODO(), 111)
	assert.NotNil(t, frame)
	assert.Equal(t, frame.(*Request).RequestHeader.CmdType, CmdTypeRequest)
	assert.Equal(t, frame.(*Request).RequestHeader.CmdCode, CmdCodeHeartbeat)
	/////// heartbeat end

	/////// response
	buf, err = bp.Encode(ctx, resp)
	assert.Nil(t, err)

	cmdInter, err = bp.Decode(ctx, buf)
	assert.Nil(t, err)
	cmdResp, ok := cmdInter.(*Response)
	assert.True(t, ok)
	assert.Equal(t, cmdResp.Content.String(), payload)
	/////// response end
}

func TestMapping(t *testing.T) {
	var (
		bp      = BoltProtocol{}
		mapping = map[uint32]uint32{
			http.StatusOK:             uint32(ResponseStatusSuccess),
			api.RouterUnavailableCode: uint32(ResponseStatusNoProcessor),
			api.NoHealthUpstreamCode:  uint32(ResponseStatusConnectionClosed),
			api.UpstreamOverFlowCode:  uint32(ResponseStatusServerThreadpoolBusy),
			api.CodecExceptionCode:    uint32(ResponseStatusCodecException),
			api.DeserialExceptionCode: uint32(ResponseStatusServerDeserialException),
			api.TimeoutExceptionCode:  uint32(ResponseStatusTimeout),
			999999:                    uint32(ResponseStatusUnknown),
		}
	)

	for k, v := range mapping {
		assert.Equal(t, bp.Mapping(k), v)
	}
}

func TestReply(t *testing.T) {
	bp := BoltProtocol{}
	// reply heartbeat
	resp := bp.Reply(context.TODO(), NewRpcResponse(1, 0, nil, buffer.NewIoBufferString("hello")))
	assert.True(t, resp.IsHeartbeatFrame())
}

func TestHijack(t *testing.T) {
	bp := BoltProtocol{}
	rsp := NewRpcResponse(1, 0, nil, buffer.NewIoBufferString("hello"))
	frame := bp.Hijack(context.TODO(), rsp, 999)
	assert.Equal(t, frame.GetStatusCode(), uint32(999))
}
