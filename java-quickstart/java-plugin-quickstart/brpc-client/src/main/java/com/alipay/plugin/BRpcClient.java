package com.alipay.plugin;

import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.Channel;
import io.netty.channel.ChannelDuplexHandler;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.util.concurrent.DefaultThreadFactory;
import java.net.InetSocketAddress;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.ReentrantLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * @author yiji@apache.org
 */
public class BRpcClient {

    Logger logger = Logger.getLogger(BRpcClient.class.getName());

    Bootstrap bootstrap;

    String host;

    int port;

    volatile Channel channel;

    static int defaultIoThreads = Math.min(Runtime.getRuntime().availableProcessors() + 1, 32);

    static final NioEventLoopGroup nioEventLoopGroup =
            new NioEventLoopGroup(defaultIoThreads, new DefaultThreadFactory("NettyClientWorker", true));

    static final Map<Integer, AsyncFuture> results = new ConcurrentHashMap<>();

    static final Map<String, BRpcClient> clients = new HashMap<>();
    static final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();
    static final ReentrantReadWriteLock.ReadLock readLock = lock.readLock();
    static final ReentrantReadWriteLock.WriteLock writeLock = lock.writeLock();

    private BRpcClient(String host, int port) {
        this.host = host;
        this.port = port;
    }

    static BRpcClient getClient(String host, int port) {
        String address = host + "_" + port;

        readLock.lock();
        try {
            BRpcClient client = clients.get(address);
            // connection is active
            if (client != null && client.channel != null && client.channel.isActive()) {
                return client;
            }
        } finally {
            readLock.unlock();
        }

        return new BRpcClient(host, port).ensureConnected();
    }

    Protocol.Response request(int id, Map<String, String> parameters, String content, int timeout) {
        Protocol.Request request = new Protocol.Request();
        request.requestId = id;

        request.timeout = timeout <= 0 ? 3000 : timeout;

        request.headers.putAll(parameters);
        request.contents = content.getBytes();

        AsyncFuture future = AsyncFuture.of(id);

        // send to remote
        ensureConnected().channel.writeAndFlush(request).addListener((ChannelFutureListener) channelFuture -> {
            if (!channelFuture.isSuccess()) {

                Protocol.Response response = new Protocol.Response();
                response.requestId = id;

                if (channelFuture.cause() != null) {
                    response.contents = channelFuture.cause().getMessage().getBytes();
                } else {
                    response.contents = ("failed to send to server: " + this.host + ":" + this.port).getBytes();
                }

                future.notify(response);
            }
        });

        return future.getResult(request.timeout);

    }

    Channel connect() throws RuntimeException {

        bootstrap = new Bootstrap();
        bootstrap.group(nioEventLoopGroup).option(ChannelOption.SO_KEEPALIVE, true)
                .option(ChannelOption.TCP_NODELAY, true).option(ChannelOption.ALLOCATOR, PooledByteBufAllocator.DEFAULT)
                .channel(NioSocketChannel.class);

        bootstrap.option(ChannelOption.CONNECT_TIMEOUT_MILLIS, 3000);

        bootstrap.handler(new ChannelInitializer() {

            @Override
            protected void initChannel(Channel ch) throws Exception {
                NettyCodecAdapter adapter = new NettyCodecAdapter();
                ch.pipeline().addLast("decoder", adapter.getDecoder()).addLast("encoder", adapter.getEncoder())
                        .addLast("handler", new NettyClientHandler());
            }
        });

        ChannelFuture future = bootstrap.connect(new InetSocketAddress(host, port));

        boolean complete = future.awaitUninterruptibly(3000, TimeUnit.MILLISECONDS);
        if (complete && future.isSuccess()) {
            Channel newChannel = future.channel();
            try {
                // Close old channel
                Channel oldChannel = BRpcClient.this.channel; // copy reference
                if (oldChannel != null) {
                    try {
                        if (logger.isLoggable(Level.INFO)) {
                            logger.info("Close old netty channel " + oldChannel + " on create new netty channel "
                                    + newChannel);
                        }
                        oldChannel.close();
                    } finally {
                    }
                }
            } finally {
                BRpcClient.this.channel = newChannel;
            }
        } else {
            throw new RuntimeException(
                    "client failed to connect to server " + host + ":" + port + ", error message is:" + (
                            future.cause() == null ? "unknown" : future.cause()));
        }

        return future.channel();
    }

    BRpcClient ensureConnected() {
        String address = this.host + "_" + this.port;

        readLock.lock();
        try {
            BRpcClient client = clients.get(address);
            // connection is active
            if (client != null && client.channel != null && client.channel.isActive()) {
                return client;
            }
        } finally {
            readLock.unlock();
        }

        writeLock.lock();
        try {

            BRpcClient client = clients.get(address);
            // double check: connection is active
            if (client != null && client.channel != null && client.channel.isActive()) {
                return client;
            }

            // client is not init or connection is closed already.
            if (client == null || client.channel == null || !client.channel.isActive()) {
                if (this.connect().isActive()) {
                    client = this;

                    // put into client cache.
                    this.clients.put(address, client);
                }
            }

            return client;
        } finally {
            writeLock.unlock();
        }

    }

    @ChannelHandler.Sharable
    static class NettyClientHandler extends ChannelDuplexHandler {

        @Override
        public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {

            if (msg instanceof Protocol.Response) {

                Protocol.Response response = (Protocol.Response) msg;
                System.out.println("receive response: " + response);

                AsyncFuture future = results.get(response.requestId);
                if (future != null) {
                    future.notify(response);
                }
            }
        }
    }

    static class AsyncFuture {

        int requestId;
        Protocol.Response response;
        ReentrantLock lock = new ReentrantLock();
        Condition notAck = lock.newCondition();

        private AsyncFuture(int id) {
            this.requestId = id;
        }

        static AsyncFuture of(int id) {
            AsyncFuture future = new AsyncFuture(id);
            BRpcClient.results.put(id, future);

            return future;
        }

        /**
         * waiting async response.
         *
         * @param milliseconds
         * @return response
         */
        Protocol.Response getResult(long milliseconds) {

            if (response != null) {
                return response;
            }

            lock.lock();
            try {
                boolean notTimeout = notAck.await(milliseconds, TimeUnit.MILLISECONDS);
                if (!notTimeout) {
                    this.response = new Protocol.Response();
                    response.requestId = requestId;
                    response.contents = "rpc invoke timeout".getBytes();
                }
            } catch (Throwable e) {
                Protocol.Response response = new Protocol.Response();
                response.requestId = requestId;
                response.contents = e.getMessage().getBytes();
                // inject failed response.
                this.response = response;
            } finally {
                lock.unlock();
            }

            return response;
        }

        void notify(Protocol.Response response) {
            lock.lock();
            try {
                this.response = response;
                notAck.signal();
            } finally {
                lock.unlock();
            }
        }
    }
}