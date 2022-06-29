package com.alipay.plugin;

import static com.alipay.plugin.Protocol.CmdRequest;
import static com.alipay.plugin.Protocol.CmdRequestHeartbeat;
import static com.alipay.plugin.Protocol.CmdResponse;
import static com.alipay.plugin.Protocol.CmdResponseHeartbeat;
import static com.alipay.plugin.Protocol.ContentLengthIndex;
import static com.alipay.plugin.Protocol.HeaderLengthIndex;
import static com.alipay.plugin.Protocol.RequestHeaderLen;
import static com.alipay.plugin.Protocol.decodeHeader;
import static com.alipay.plugin.Protocol.encodeHeader;
import static com.alipay.plugin.Protocol.getHeaderLength;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.ByteToMessageDecoder;
import io.netty.handler.codec.MessageToByteEncoder;
import java.util.List;

/**
 * @author yiji@apache.org
 */
public class NettyCodecAdapter {

    private final ChannelHandler encoder = new InternalEncoder();

    private final ChannelHandler decoder = new InternalDecoder();

    public ChannelHandler getEncoder() {
        return encoder;
    }

    public ChannelHandler getDecoder() {
        return decoder;
    }

    private class InternalEncoder extends MessageToByteEncoder {

        @Override
        protected void encode(ChannelHandlerContext ctx, Object msg, ByteBuf out) throws Exception {
            if (msg instanceof Protocol.Request) {

                Protocol.Request request = (Protocol.Request) msg;
                if (request.headers != null && request.headers.size() > 0) {
                    request.headerLength = getHeaderLength(request.headers);
                }
                if (request.contents != null && request.contents.length > 0) {
                    request.contentLength = (short) request.contents.length;
                }

                out.writeShort(request.magic);
                out.writeByte(request.flag);
                out.writeInt(request.requestId);
                out.writeByte(request.codec);
                out.writeInt(request.timeout);
                out.writeShort(request.headerLength);
                out.writeShort(request.contentLength);

                if (request.headerLength > 0) {
                    encodeHeader(out, request.headers);
                }

                if (request.contentLength > 0) {
                    out.writeBytes(request.contents);
                }

            } else if (msg instanceof Protocol.Response) {

                Protocol.Response response = (Protocol.Response) msg;

                if (response.headers != null && response.headers.size() > 0) {
                    response.headerLength = getHeaderLength(response.headers);
                }
                if (response.contents != null && response.contents.length > 0) {
                    response.contentLength = (short) response.contents.length;
                }

                out.writeShort(response.magic);
                out.writeByte(response.flag);
                out.writeInt(response.requestId);
                out.writeByte(response.codec);
                out.writeInt(response.status);
                out.writeShort(response.headerLength);
                out.writeShort(response.contentLength);

                if (response.headerLength > 0) {
                    encodeHeader(out, response.headers);
                }

                if (response.contentLength > 0) {
                    out.writeBytes(response.contents);
                }
            }
        }
    }

    private class InternalDecoder extends ByteToMessageDecoder {

        @Override
        protected void decode(ChannelHandlerContext ctx, ByteBuf input, List<Object> out) throws Exception {
            // decode object.

            do {
                int savedReaderIndex = input.readerIndex();

                if (input.readableBytes() < RequestHeaderLen) {
                    break;
                }

                short headerLength = input.getShort(savedReaderIndex + HeaderLengthIndex);
                short contentLength = input.getShort(savedReaderIndex + ContentLengthIndex);
                int frameLength = RequestHeaderLen + headerLength + contentLength;
                if (input.readableBytes() < frameLength) {
                    break;
                }

                byte flag = input.getByte(savedReaderIndex + 2);
                switch (flag) {
                    case CmdRequest:
                    case CmdRequestHeartbeat:
                        Protocol.Request request = new Protocol.Request();
                        request.magic = input.readShort();
                        request.flag = input.readByte();
                        request.requestId = input.readInt();
                        request.codec = input.readByte();
                        request.timeout = input.readInt();
                        request.headerLength = input.readShort();
                        request.contentLength = input.readShort();

                        if (request.headerLength > 0) {
                            byte[] headers = new byte[request.headerLength];
                            input.readBytes(headers);
                            decodeHeader(headers, request.headers);
                        }

                        if (request.contentLength > 0) {
                            byte[] contents = new byte[request.contentLength];
                            input.readBytes(contents);
                            request.contents = contents;
                        }

                        out.add(request);

                        break;
                    case CmdResponse:
                    case CmdResponseHeartbeat:
                        Protocol.Response response = new Protocol.Response();
                        response.magic = input.readShort();
                        response.flag = input.readByte();
                        response.requestId = input.readInt();
                        response.codec = input.readByte();
                        response.status = input.readInt();
                        response.headerLength = input.readShort();
                        response.contentLength = input.readShort();

                        if (response.headerLength > 0) {
                            byte[] headers = new byte[response.headerLength];
                            input.readBytes(headers);
                            decodeHeader(headers, response.headers);
                        }

                        if (response.contentLength > 0) {
                            byte[] contents = new byte[response.contentLength];
                            input.readBytes(contents);
                            response.contents = contents;
                        }

                        out.add(response);
                }

            } while (input.readableBytes() > 0);
        }
    }

}