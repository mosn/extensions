package com.alipay.plugin;

import static com.alipay.plugin.XRpcServerBootStrap.port;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.Channel;
import io.netty.channel.ChannelDuplexHandler;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.codec.ByteToMessageDecoder;
import io.netty.handler.codec.MessageToByteEncoder;
import io.netty.util.concurrent.DefaultThreadFactory;
import java.util.List;
import java.util.concurrent.CountDownLatch;

/**
 * @author yiji@apache.org
 */
public class XRpcController {

    public static void main(String[] args) throws InterruptedException {
        startSocketService(port);
    }

    public static void startSocketService(int port) throws InterruptedException {

        CountDownLatch latch = new CountDownLatch(1);

        ServerBootstrap bootstrap = new ServerBootstrap();
        NioEventLoopGroup bossGroup = new NioEventLoopGroup(1, new DefaultThreadFactory("NettyServerBoss", true));
        NioEventLoopGroup workerGroup = new NioEventLoopGroup(2, new DefaultThreadFactory("NettyServerWorker", true));

        bootstrap.group(bossGroup, workerGroup)
                .channel(NioServerSocketChannel.class)
                .childOption(ChannelOption.TCP_NODELAY, Boolean.TRUE)
                .childOption(ChannelOption.SO_REUSEADDR, Boolean.TRUE)
                .childOption(ChannelOption.ALLOCATOR, PooledByteBufAllocator.DEFAULT)
                .childHandler(new ChannelInitializer<NioSocketChannel>() {
                    @Override
                    protected void initChannel(NioSocketChannel ch) throws Exception {
                        NettyCodecAdapter adapter = new NettyCodecAdapter();
                        ch.pipeline()
                                .addLast("decoder", adapter.getDecoder())
                                .addLast("encoder", adapter.getEncoder())
                                .addLast("handler", new NettyServerHandler());
                    }
                });
        // bind
        ChannelFuture channelFuture = bootstrap.bind("0.0.0.0", port);
        channelFuture.syncUninterruptibly();
        Channel channel = channelFuture.channel();

        latch.await();
    }

    static class NettyCodecAdapter {

        static final int HEADER_LENGTH = 10;

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
                // append prefix of '0'
                String packet = String.valueOf(msg);
                String lenOfPacket = String.valueOf(packet.getBytes().length);
                if (lenOfPacket.length() < HEADER_LENGTH) {
                    int remain = HEADER_LENGTH - lenOfPacket.length();
                    String prefix = "";
                    for (int i = 0; i < remain; i++) {
                        prefix += "0";
                    }
                    lenOfPacket = prefix + lenOfPacket;
                }
                // write length + body
                out.writeBytes(lenOfPacket.getBytes());
                out.writeBytes(packet.getBytes());
            }
        }

        private class InternalDecoder extends ByteToMessageDecoder {

            @Override
            protected void decode(ChannelHandlerContext ctx, ByteBuf input, List<Object> out) throws Exception {
                // decode object.
                do {
                    int saveReaderIndex = input.readerIndex();
                    // check min readable protocol length
                    if (input.readableBytes() < HEADER_LENGTH) {
                        break;
                    }

                    byte[] header = new byte[HEADER_LENGTH];
                    input.readBytes(header);

                    // resolve protocol header length
                    String length = new String(header);
                    while (length.startsWith("0")) {
                        if (length.length() > 1) {
                            length = length.substring(1);
                        }
                    }

                    int packetLength = Integer.parseInt(length);
                    int available = input.readableBytes();
                    // check protocol full package length
                    if (packetLength > available) {
                        // rollback buffer reader index
                        input.readerIndex(saveReaderIndex);
                        break;
                    }

                    byte[] payload = new byte[packetLength];
                    input.readBytes(payload);

                    // decode full request or response
                    out.add(new String(payload));

                    // skip telnet \r\n if exists.
                    if (input.readableBytes() > 0) {
                        byte b = input.readByte();
                        while (b == '\r' || b == '\n') {
                            if (input.isReadable()) {
                                b = input.readByte();
                            } else {
                                break;
                            }
                        }

                        // rollback read 1 byte
                        if (b != '\r' && b != '\n') {
                            input.readerIndex(input.readerIndex() - 1);
                        }
                    }

                } while (input.readableBytes() > 0);
            }
        }
    }

    @io.netty.channel.ChannelHandler.Sharable
    static class NettyServerHandler extends ChannelDuplexHandler {
        @Override
        public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {

            String body = String.valueOf(msg);

            body = body.replaceAll("<RequestFlag>0</RequestFlag>", "<RequestFlag>1</RequestFlag>");

            int index = body.indexOf("</Body>");

            StringBuilder sb = new StringBuilder();

            if (index >= 0) {
                sb.append(body, 0, index);

                if (body.contains("<userId>")) {
                    sb.append("  <title>developer</title>\n");
                    sb.append("    <address>hangzhou</address>\n");
                }

                sb.append("    <response_ip>")
                        .append(ctx.channel().localAddress().toString().substring(1))
                        .append("</response_ip>\n");

                sb.append(body.substring(index));
            } else {
                sb.append(body);
            }

            // inject server ip:port

            System.out.println("reply response: " + sb.toString());

            ctx.channel().writeAndFlush(sb.toString());
        }
    }

}