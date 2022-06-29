package com.alipay.plugin;

import static com.alipay.plugin.BRpcServerBootstrap.port;
import static com.alipay.plugin.Protocol.CmdRequestHeartbeat;
import static com.alipay.plugin.Protocol.CmdResponseHeartbeat;
import io.netty.bootstrap.ServerBootstrap;
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
import io.netty.util.concurrent.DefaultThreadFactory;
import java.util.concurrent.CountDownLatch;

/**
 * @author yiji@apache.org
 */
public class BRpcController {

    public static void main(String[] args) throws InterruptedException {
        startSocketService(port);
    }


    public static void startSocketService(int port) throws InterruptedException {

        CountDownLatch latch = new CountDownLatch(1);

        ServerBootstrap bootstrap = new ServerBootstrap();
        NioEventLoopGroup bossGroup = new NioEventLoopGroup(1, new DefaultThreadFactory("NettyServerBoss", true));
        NioEventLoopGroup workerGroup = new NioEventLoopGroup(2, new DefaultThreadFactory("NettyServerWorker", true));

        bootstrap.group(bossGroup, workerGroup).channel(NioServerSocketChannel.class)
                .childOption(ChannelOption.TCP_NODELAY, Boolean.TRUE)
                .childOption(ChannelOption.SO_REUSEADDR, Boolean.TRUE)
                .childOption(ChannelOption.ALLOCATOR, PooledByteBufAllocator.DEFAULT)
                .childHandler(new ChannelInitializer<NioSocketChannel>() {
                    @Override
                    protected void initChannel(NioSocketChannel ch) throws Exception {
                        NettyCodecAdapter adapter = new NettyCodecAdapter();
                        ch.pipeline().addLast("decoder", adapter.getDecoder()).addLast("encoder", adapter.getEncoder())
                                .addLast("handler", new NettyServerHandler());
                    }
                });
        // bind
        ChannelFuture channelFuture = bootstrap.bind("0.0.0.0", port);
        channelFuture.syncUninterruptibly();
        Channel channel = channelFuture.channel();
        latch.await();
    }

    @ChannelHandler.Sharable
    static class NettyServerHandler extends ChannelDuplexHandler {
        @Override
        public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {

            if (msg instanceof Protocol.Request) {

                Protocol.Request request = (Protocol.Request) msg;
                System.out.println("receive request: " + request);

                Protocol.Response response = new Protocol.Response();
                response.requestId = request.requestId;

                if (request.headers != null && request.headers.size() > 0) {
                    response.headers.putAll(request.headers);
                }

                // inject who I am.
                try {
                    response.headers.put("response_ip", ctx.channel().localAddress().toString().substring(1));
                } catch (Throwable ignored) {
                }

                if (request.contents != null && request.contents.length > 0) {
                    response.contentLength = (short) request.contents.length;
                    response.contents = request.contents;
                }

                if (request.flag == CmdRequestHeartbeat) {
                    response.flag = CmdResponseHeartbeat;
                    response.headers.put("response_heartbeat", "yes");
                }

                System.out.println("reply response: " + response);

                ctx.channel().writeAndFlush(response);
            } else if (msg instanceof Protocol.Response) {
                Protocol.Response response = (Protocol.Response) msg;
                System.out.println("receive response: " + response);

            }
        }
    }
}
