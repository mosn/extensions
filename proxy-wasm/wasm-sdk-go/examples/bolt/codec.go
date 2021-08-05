package bolt

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy"
	"reflect"
)

type boltCodec struct {
}

func (c *boltCodec) Encode(ctx context.Context, cmd proxy.Command) (proxy.Buffer, error) {
	// encode request command
	if req, ok := cmd.(*Request); ok {
		return encodeRequest(ctx, req)
	}

	// encode response command
	if resp, ok := cmd.(*Response); ok {
		return encodeResponse(ctx, resp)
	}

	proxy.Log.Warnf("[boltCodec] maybe receive unsupported command type: %v", reflect.TypeOf(cmd))
	return nil, fmt.Errorf("unknow command type: %v", reflect.TypeOf(cmd))
}

func (c *boltCodec) Decode(ctx context.Context, data proxy.Buffer) (proxy.Command, error) {
	if data.Len() >= LessLen {
		cmdType := data.Bytes()[1]

		switch cmdType {
		case CmdTypeRequest:
			return decodeRequest(ctx, data, false)
		case CmdTypeRequestOneway:
			return decodeRequest(ctx, data, true)
		case CmdTypeResponse:
			return decodeResponse(ctx, data)
		default:
			// unknown cmd type
			return nil, fmt.Errorf("decode Error, type = %s, value = %d", UnKnownCmdType, cmdType)
		}
	}

	return nil, nil
}

func decodeRequest(ctx context.Context, data proxy.Buffer, oneway bool) (cmd proxy.Command, err error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is RequestHeaderLen(22)
	if bytesLen < RequestHeaderLen {
		return
	}

	// 2. least bytes to decode whole frame
	classLen := binary.BigEndian.Uint16(bytes[14:16])
	headerLen := binary.BigEndian.Uint16(bytes[16:18])
	contentLen := binary.BigEndian.Uint32(bytes[18:22])

	frameLen := RequestHeaderLen + int(classLen) + int(headerLen) + int(contentLen)
	if bytesLen < frameLen {
		return
	}
	data.Drain(frameLen)

	// 3. decode header
	request := &Request{}

	cmdType := CmdTypeRequest
	if oneway {
		cmdType = CmdTypeRequestOneway
	}

	request.RpcHeader = RpcHeader{
		Protocol:   ProtocolCode,
		CmdType:    cmdType,
		CmdCode:    binary.BigEndian.Uint16(bytes[2:4]),
		Version:    bytes[4],
		RequestId:  binary.BigEndian.Uint32(bytes[5:9]),
		Codec:      bytes[9],
		ClassLen:   classLen,
		HeaderLen:  headerLen,
		ContentLen: contentLen,
	}
	request.Timeout = int32(binary.BigEndian.Uint32(bytes[10:14]))

	request.Data = proxy.NewBuffer(frameLen)

	//5. copy data for io multiplexing
	request.Data.Write(bytes[:frameLen])
	request.rawData = request.Data.Bytes()

	//6. process wrappers: Class, Header, Content, Data
	headerIndex := RequestHeaderLen + int(classLen)
	contentIndex := headerIndex + int(headerLen)

	request.rawMeta = request.rawData[:RequestHeaderLen]
	if classLen > 0 {
		request.rawClass = request.rawData[RequestHeaderLen:headerIndex]
		request.Class = string(request.rawClass)
	}
	if headerLen > 0 {
		request.rawHeader = request.rawData[headerIndex:contentIndex]
		err = proxy.DecodeHeader(request.rawHeader, &request.CommonHeader)
	}
	if contentLen > 0 {
		request.rawContent = request.rawData[contentIndex:]
		request.Content = proxy.WrapBuffer(request.rawContent)
	}
	return request, err
}

func decodeResponse(ctx context.Context, data proxy.Buffer) (cmd proxy.Command, err error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is ResponseHeaderLen(20)
	if bytesLen < ResponseHeaderLen {
		return
	}

	// 2. least bytes to decode whole frame
	classLen := binary.BigEndian.Uint16(bytes[12:14])
	headerLen := binary.BigEndian.Uint16(bytes[14:16])
	contentLen := binary.BigEndian.Uint32(bytes[16:20])

	frameLen := ResponseHeaderLen + int(classLen) + int(headerLen) + int(contentLen)
	if bytesLen < frameLen {
		return
	}
	data.Drain(frameLen)

	// 3. decode header
	response := &Response{}

	response.RpcHeader = RpcHeader{
		Protocol:   ProtocolCode,
		CmdType:    CmdTypeResponse,
		CmdCode:    binary.BigEndian.Uint16(bytes[2:4]),
		Version:    bytes[4],
		RequestId:  binary.BigEndian.Uint32(bytes[5:9]),
		Codec:      bytes[9],
		ClassLen:   classLen,
		HeaderLen:  headerLen,
		ContentLen: contentLen,
	}
	response.Status = binary.BigEndian.Uint16(bytes[10:12])

	response.Data = proxy.NewBuffer(frameLen)

	//TODO: test recycle by model, so we can recycle request/response models, headers also
	//4. copy data for io multiplexing
	response.Data.Write(bytes[:frameLen])
	response.rawData = response.Data.Bytes()

	//5. process wrappers: Class, Header, Content, Data
	headerIndex := ResponseHeaderLen + int(classLen)
	contentIndex := headerIndex + int(headerLen)

	response.rawMeta = response.rawData[:ResponseHeaderLen]
	if classLen > 0 {
		response.rawClass = response.rawData[ResponseHeaderLen:headerIndex]
		response.Class = string(response.rawClass)
	}
	if headerLen > 0 {
		response.rawHeader = response.rawData[headerIndex:contentIndex]
		err = proxy.DecodeHeader(response.rawHeader, &response.CommonHeader)
	}
	if contentLen > 0 {
		response.rawContent = response.rawData[contentIndex:]
		response.Content = proxy.WrapBuffer(response.rawContent)
	}
	return response, err
}

func encodeRequest(ctx context.Context, request *Request) (proxy.Buffer, error) {
	// 1. fast-path, use existed raw data
	if request.rawData != nil {
		// 1.1 replace requestId
		binary.BigEndian.PutUint32(request.rawMeta[RequestIdIndex:], request.RequestId)

		// 1.2 check if header/content changed
		if !request.RpcHeader.Changed && !request.ContentChanged {
			return request.Data, nil
		}
	}

	// 2. slow-path, construct buffer from scratch

	// 2.1 calculate frame length
	if request.Class != "" {
		request.ClassLen = uint16(len(request.Class))
	}
	if request.CommonHeader.Size() != 0 {
		request.HeaderLen = uint16(proxy.GetEncodeHeaderLength(&request.CommonHeader))
	}
	if request.Content != nil {
		request.ContentLen = uint32(request.Content.Len())
	}
	frameLen := RequestHeaderLen + int(request.ClassLen) + int(request.HeaderLen) + int(request.ContentLen)

	// 2.2 alloc encode buffer, this buffer will be recycled after connection.Write
	buf := proxy.NewBuffer(frameLen)

	// 2.3 encode: meta, class, header, content
	// 2.3.1 meta
	buf.WriteByte(request.Protocol)
	buf.WriteByte(request.CmdType)
	buf.WriteUint16(request.CmdCode)
	buf.WriteByte(request.Version)
	buf.WriteUint32(request.RequestId)
	buf.WriteByte(request.Codec)
	buf.WriteUint32(uint32(request.Timeout))
	buf.WriteUint16(request.ClassLen)
	buf.WriteUint16(request.HeaderLen)
	buf.WriteUint32(request.ContentLen)
	// 2.3.2 class
	if request.ClassLen > 0 {
		buf.WriteString(request.Class)
	}
	// 2.3.3 header
	if request.HeaderLen > 0 {
		proxy.EncodeHeader(buf, &request.CommonHeader)
	}
	// 2.3.4 content
	if request.ContentLen > 0 {
		// use request.Content.WriteTo might have error under retry scene
		buf.Write(request.Content.Bytes())
	}

	return buf, nil
}

func encodeResponse(ctx context.Context, response *Response) (proxy.Buffer, error) {
	// 1. fast-path, use existed raw data
	if response.rawData != nil {
		// 1. replace requestId
		binary.BigEndian.PutUint32(response.rawMeta[RequestIdIndex:], uint32(response.RequestId))

		// 2. check header change
		if !response.CommonHeader.Changed && !response.ContentChanged {
			return response.Data, nil
		}
	}

	// 2. slow-path, construct buffer from scratch

	// 2.1 calculate frame length
	if response.Class != "" {
		response.ClassLen = uint16(len(response.Class))
	}
	if response.CommonHeader.Size() != 0 {
		response.HeaderLen = uint16(proxy.GetEncodeHeaderLength(&response.CommonHeader))
	}
	if response.Content != nil {
		response.ContentLen = uint32(response.Content.Len())
	}
	frameLen := ResponseHeaderLen + int(response.ClassLen) + int(response.HeaderLen) + int(response.ContentLen)

	// 2.2 alloc encode buffer, this buffer will be recycled after connection.Write
	buf := proxy.NewBuffer(frameLen)

	// 2.3 encode: meta, class, header, content
	// 2.3.1 meta
	buf.WriteByte(response.Protocol)
	buf.WriteByte(response.CmdType)
	buf.WriteUint16(response.CmdCode)
	buf.WriteByte(response.Version)
	buf.WriteUint32(response.RequestId)
	buf.WriteByte(response.Codec)
	buf.WriteUint16(response.Status)
	buf.WriteUint16(response.ClassLen)
	buf.WriteUint16(response.HeaderLen)
	buf.WriteUint32(response.ContentLen)
	// 2.3.2 class
	if response.ClassLen > 0 {
		buf.WriteString(response.Class)
	}
	// 2.3.3 header
	if response.HeaderLen > 0 {
		proxy.EncodeHeader(buf, &response.CommonHeader)
	}
	// 2.3.4 content
	if response.ContentLen > 0 {
		// use request.Content.WriteTo might have error under retry scene
		buf.Write(response.Content.Bytes())
	}

	return buf, nil
}
