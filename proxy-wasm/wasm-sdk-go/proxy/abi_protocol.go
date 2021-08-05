package proxy

import (
	"context"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
)

//export proxy_decode_buffer_bytes
func proxyDecodeBufferBytes(contextID uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode buffer, contextId %d not found", contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of bytes to be parsed
	data := parseByteSlice(bufferData, len)
	buffer := WrapBuffer(data)
	// call user extension implementation
	cmd, err := ctx.Codec().Decode(context.TODO(), buffer)
	if err != nil {
		log.Fatalf("failed to decode buffer by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	// need more data
	if cmd == nil {
		return types.StatusNeedMoreData
	}

	if buffer.Pos() == 0 {
		// When decoding is complete, the contents of the buffer should be read
		log.Errorf("the contents of the buffer should be read by protocol %s, contextId %v, buffer pos %v", ctx.Name(), contextID, buffer.Pos())
		return types.StatusInternalFailure
	}

	ctx.(attribute).set(types.AttributeKeyDecodeCommand, cmd)

	decode := decodeCommandBuffer(cmd, buffer.Pos(), contextID)
	// report encode data
	err = setDecodeBuffer(decode.Bytes())
	if err != nil {
		log.Errorf("failed to report decode buffer by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	return types.StatusOK
}

//export proxy_encode_buffer_bytes
func proxyEncodeBufferBytes(contextID uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to encode buffer, contextId %v not found", contextID)
		return types.StatusInternalFailure
	}
	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of dataBytes to be parsed
	data := parseByteSlice(bufferData, len)
	buffer := WrapBuffer(data)

	// bufferData format:
	// encoded header map | Flag | replaceId, id | (Timeout|Status) | drain length | raw dataBytes
	headerBytes, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to read decode buffer header map, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	headers := &CommonHeader{}
	// encoded header map
	if headerBytes > 0 {
		DecodeHeader(data[4:4+headerBytes], headers)
	}
	// skip header bytes
	buffer.Drain(headerBytes)

	flag, err := buffer.ReadByte()
	if err != nil {
		log.Errorf("failed to decode buffer flag, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	attr := ctx.(attribute)

	// find context cmd
	cachedCmd := attr.attr(types.AttributeKeyDecodeCommand)

	if cmd := attr.attr(types.AttributeKeyEncodeCommand); cmd != nil {
		// reply heartbeat ?
		if parsedCmd, ok := cmd.(Command); ok && parsedCmd.IsHeartbeat() {
			cachedCmd = cmd
		}
	}

	if cachedCmd == nil {
		// is heartbeat 、keep-alive or hijack ?
		cachedCmd = attr.attr(types.AttributeKeyEncodeCommand)
	}

	if cachedCmd == nil {
		log.Errorf("failed to find cached command, maybe a bug occurred, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	// Multiplexing ID: This is equivalent to the stream ID
	replacedId, err := buffer.ReadUint64()
	if err != nil {
		log.Errorf("failed to decode buffer replacedId, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	var cmd Command
	cmdType := flag >> 6
	switch cmdType {
	case types.RequestType,
		types.RequestOneWayType:
		cmd, ok = cachedCmd.(Request)
		if !ok {
			detect, _ := cachedCmd.(Command)
			log.Errorf("cached cmd should be Request, maybe a bug occurred, contextId: %v, actual hb: %v", contextID, detect.IsHeartbeat())
			return types.StatusInternalFailure
		}

	case types.ResponseType:
		cmd, ok = cachedCmd.(Response)
		if !ok {
			detect, _ := cachedCmd.(Command)
			log.Errorf("cached cmd should be Response, maybe a bug occurred, contextId: %v, actual hb: %v ", contextID, detect.IsHeartbeat())
			return types.StatusInternalFailure
		}
	default:
		log.Errorf("failed to decode buffer, type = %s, value = %d", types.UnKnownRpcFlagType, flag)
		return types.StatusInternalFailure
	}

	id, err := buffer.ReadUint64()
	// we check encoded id equals cached command id
	if id != cmd.CommandId() {
		log.Errorf("encode buffer command id is not match, contextId: %v , cached id = %d, actual = %d, replaced id = %d", contextID, cmd.CommandId(), id, replacedId)
		return types.StatusInternalFailure
	}

	// skip timeout or status
	buffer.ReadInt()

	dataBytes, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to decode buffer drain length, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	if dataBytes > 0 {
		cmd.SetData(WrapBuffer(data[buffer.Pos():]))
	}

	// override cached request
	injectHeaderIfRequired(cmd, headers)

	// update command replacedId
	cmd.SetCommandId(replacedId)
	// call user extension implementation
	encode, err := ctx.Codec().Encode(context.TODO(), cmd)
	if err != nil {
		log.Fatalf("failed to encode command by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
		return types.StatusInternalFailure
	}

	attr.set(types.AttributeKeyEncodedBuffer, encode)

	// we don't need encode header again, the host side only pays attention to
	// the buffer of encode and sends it directly to the remote host
	proxyBuffer := encodeCommandBuffer(cmd, encode)
	// report encode data
	err = setEncodeBuffer(proxyBuffer.Bytes())
	if err != nil {
		log.Errorf("failed to report encode buffer by protocol %s, contextId %v, err %v", ctx.Name(), contextID, err)
	}

	return types.StatusOK
}

//export proxy_keepalive_buffer_bytes
func proxyKeepAliveBufferBytes(contextID uint32, id int64) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode keepalive buffer, contextId %v not found", contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	// not support keepalive
	keepAlive := ctx.KeepAlive()
	if keepAlive == nil {
		return types.StatusBadArgument
	}

	cmd := keepAlive.KeepAlive(uint64(id))
	if cmd == nil {
		return types.StatusBadArgument
	}

	attr := ctx.(attribute)
	attr.set(types.AttributeKeyEncodeCommand, cmd)

	return types.StatusOK
}

//export proxy_reply_keepalive_buffer_bytes
func proxyReplyKeepAliveBufferBytes(contextID uint32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode reply keepalive buffer, contextId %v not found", contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	cmd := ctx.(attribute).attr(types.AttributeKeyDecodeCommand)

	resp := ctx.KeepAlive().ReplyKeepAlive(cmd.(Request))
	attr := ctx.(attribute)
	attr.set(types.AttributeKeyEncodeCommand, resp)

	return types.StatusOK
}

//export proxy_hijack_buffer_bytes
func proxyHijackBufferBytes(contextID uint32, statusCode int32, bufferData *byte, len int) types.Status {
	ctx, ok := this.protocolStreams[contextID]
	if !ok {
		log.Errorf("failed to decode hijack buffer, contextId %v not found", contextID)
		return types.StatusInternalFailure
	}

	this.setActiveContextID(contextID)

	if len <= 0 {
		// should never be happen
		return types.StatusEmpty
	}

	// convert data into an array of dataBytes to be parsed
	data := parseByteSlice(bufferData, len)
	buffer := WrapBuffer(data)

	// bufferData format:
	// encoded header map | Flag | replaceId, id | Timeout | drain length | raw dataBytes
	headerBytes, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to read hijack buffer header map, contextId: %v, err: %v", contextID, err)
		return types.StatusInternalFailure
	}

	offset := 4 + headerBytes + 1 + 8*2 + 4
	// move drain length offset
	buffer.Move(offset)

	drainLen, err := buffer.ReadInt()
	if err != nil {
		log.Errorf("failed to read hijack buffer, contextId: %v, err: %v", contextID, err)
		return types.StatusInternalFailure
	}

	if drainLen <= 0 {
		log.Errorf("hijack request content is nil, contextId: %v", contextID)
		return types.StatusInternalFailure
	}

	payload := data[offset+4 : offset+4+drainLen]
	// build request command
	cmd, err := ctx.Codec().Decode(context.TODO(), WrapBuffer(payload))
	if err != nil {
		log.Errorf("failed to build hijack request command, contextId: %v, err: %v", contextID, err)
		return types.StatusInternalFailure
	}

	resp := ctx.Hijacker().Hijack(cmd.(Request), Mapping(statusCode))

	attr := ctx.(attribute)
	attr.set(types.AttributeKeyEncodeCommand, resp)

	return types.StatusOK
}

func decodeCommandBuffer(cmd Command, drainBytes int, contextID uint32) Buffer {
	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|GetStatus) | drain length | raw bytes length | raw bytes
	headers := cmd.GetHeader()
	buf := AllocateBuffer()

	headerBytes := GetEncodeHeaderLength(headers)
	buf.WriteInt(headerBytes)
	// encoded header map
	if headerBytes > 0 {
		EncodeHeader(buf, headers)
	}

	var flag byte
	if cmd.IsHeartbeat() {
		flag = HeartBeatFlag
	}

	// should copy raw bytes
	flag = flag | CopyRawBytesFlag

	// record flag write index
	flagIndex := buf.Len()
	// write flag
	buf.WriteByte(flag)
	// write id
	buf.WriteUint64(cmd.CommandId())

	// check is request
	if req, ok := cmd.(Request); ok {
		flag = flag | RpcRequestFlag
		if req.IsOneWay() {
			flag = flag | RpcOnewayFlag
		}
		// update request flag
		buf.PutByte(flagIndex, flag)
		buf.WriteUint32(req.GetTimeout())
	} else if resp, ok := cmd.(Response); ok {
		buf.WriteUint32(resp.GetStatus())
	}

	buf.WriteInt(drainBytes)
	if drainBytes > 0 {
		contentBytes := 0
		if cmd.GetData() != nil {
			contentBytes = cmd.GetData().Len()
		}
		// write decode content length
		buf.WriteInt(contentBytes)
		// write decode content, protocol header is not included
		if contentBytes > 0 {
			buf.Write(cmd.GetData().Bytes())
		}
	}

	return buf
}

func encodeCommandBuffer(cmd Command, encode Buffer) Buffer {
	// bufferData format:
	// encoded header map | Flag | Id | (Timeout|GetStatus) | drain length | raw bytes
	buf := AllocateBuffer()

	var headerBytes = 0
	buf.WriteInt(headerBytes)

	var flag byte
	if cmd.IsHeartbeat() {
		flag = HeartBeatFlag
	}

	// should copy raw bytes
	flag = flag | CopyRawBytesFlag

	// record flag index
	flagIndex := buf.Len()
	// write flag
	buf.WriteByte(flag)
	// write id
	buf.WriteUint64(cmd.CommandId())

	// check is request
	if req, ok := cmd.(Request); ok {
		flag = flag | RpcRequestFlag
		if req.IsOneWay() {
			flag = flag | RpcOnewayFlag
		}
		// update request flag
		buf.PutByte(flagIndex, flag)
		buf.WriteUint32(req.GetTimeout())
	} else if resp, ok := cmd.(Response); ok {
		buf.WriteUint32(resp.GetStatus())
	}

	drainBytes := encode.Len()
	buf.WriteInt(drainBytes)
	if drainBytes > 0 {
		buf.Write(encode.Bytes())
	}

	return buf
}

func injectHeaderIfRequired(cmd Command, headers *CommonHeader) {
	if cmd.GetHeader().Size() != headers.Size() {
		cmd.GetHeader().Range(func(key, value string) bool {
			v, ok := headers.Get(key)
			if !ok {
				// remove old key
				cmd.GetHeader().Del(key)
			} else {
				// add new key
				cmd.GetHeader().Set(key, v)
			}
			return true
		})
	}
}
