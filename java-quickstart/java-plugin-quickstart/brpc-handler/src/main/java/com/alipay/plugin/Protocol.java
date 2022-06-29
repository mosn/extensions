package com.alipay.plugin;

import io.netty.buffer.ByteBuf;
import java.util.HashMap;
import java.util.Map;

/**
 * @author yiji@apache.org
 *
 * 0    1    2      3   4       6         8        10         12          14          16
 * +----+----+------+---+---+---+---+-----+----+----+----+----+-----+-----+------+-----+
 * |  magic  | flag |   requestID   |codec|    timeout/status | headerLen | contentLen |
 * +----+----+------+---+---+---+---+-----+----+----+----+----+-----+-----+------+-----+
 * |  header   + content  bytes     ... ...                                            |
 * |                                                                                   |
 * +-----------------------------------------------------------------------------------+
 *
 * 字段说明：
 * ● 2字节 magic 魔法数，固定值： 0xbcbc
 * ● 1字节 flag 报文标志，取值 1：请求， 2：响应，3： 心跳请求，4：心跳响应
 * ● 4字节 requestID 请求或者响应id
 * ● 1字节 codec 序列化编号，固定值：0
 * ● 4字节 timeout 超时时间（flag取值1、3时有效，代表请求超时），请求时该字段填值
 * ● 4字节 status 响应码（flag取值2、4时有效，代表响应状态码），响应时该字段填值，和timeout公用字段
 * ● 2字节 headerLen代表报文key value键值对长度
 * ● 2字节 contentLen 代表消息体长度
 * ● header 编码后的格式为键值对字符串，格式如下：
 * ○ key1=value1&key2=value2   举例：interface=com.alipay.core.UserService&method=userInfo
 * ● content bytes 为业务报文体，内容不设限制，比如传递字符串
 */
public class Protocol {

    static final int RequestHeaderLen = 16;
    static final int LessLen = RequestHeaderLen;

    static final int HeaderLengthIndex = 12;
    static final int ContentLengthIndex = HeaderLengthIndex + 2;


    static final byte CmdRequest = 1;
    static final byte CmdResponse = 2;
    static final byte CmdRequestHeartbeat = 3;
    static final byte CmdResponseHeartbeat = 4;

    public static class Request {

        short magic = (short) 0xbcbc; // fixed 0xbcbc
        byte flag = CmdRequest;
        int requestId;
        byte codec;
        int timeout;
        short headerLength;
        short contentLength;

        Map<String, String> headers = new HashMap<>();
        byte[] contents;

        @Override
        public String toString() {
            return "Request{" +
                    "magic=" + magic +
                    ", flag=" + flag + (flag == CmdRequestHeartbeat ? " heartbeat" : "") +
                    ", requestId=" + requestId +
                    ", codec=" + codec +
                    ", timeout=" + timeout +
                    ", headerLength=" + headerLength +
                    ", contentLength=" + contentLength +
                    ", headers=" + headers +
                    ", contents=" + ((contents == null) ? "" : new String(contents)) +
                    '}';
        }
    }

    public static class Response {

        short magic = (short) 0xbcbc;
        byte flag = CmdResponse;
        int requestId;
        byte codec;
        int status;
        short headerLength;
        short contentLength;

        Map<String, String> headers = new HashMap<>();
        byte[] contents;

        @Override
        public String toString() {
            return "Response{" +
                    "magic=" + magic +
                    ", flag=" + flag + (flag == CmdResponseHeartbeat ? " heartbeat" : "") +
                    ", requestId=" + requestId +
                    ", codec=" + codec +
                    ", status=" + status +
                    ", headerLength=" + headerLength +
                    ", contentLength=" + contentLength +
                    ", headers=" + headers +
                    ", contents=" + ((contents == null) ? "" : new String(contents)) +
                    '}';
        }
    }

    static void decodeHeader(byte[] bytes, Map<String, String> h) {
        if (bytes.length > 0) {
            String str = new String(bytes);
            String[] items = str.split("&");
            for (String item : items) {
                String[] pair = item.split("=");
                if (pair.length == 2) {
                    h.put(pair[0], pair[1]);
                }
            }
        }
    }

    static int encodeHeader(ByteBuf buf, Map<String, String> h) {
        StringBuffer sb = new StringBuffer();

        for (String key : h.keySet()) {
            if (sb.length() > 0) {
                sb.append("&");
            }

            sb.append(key);
            sb.append("=");
            sb.append(h.get(key));
        }

        byte[] bytes = sb.toString().getBytes();
        if (buf != null) {
            buf.writeBytes(bytes);
        }

        return bytes.length;
    }

    static short getHeaderLength(Map<String, String> h) {
        short total = (short) encodeHeader(null, h);
        return total;
    }

}