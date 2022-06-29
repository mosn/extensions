package com.alipay.plugin;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yiji@apache.org
 */
@RestController()
public class InvokeController {

    @GetMapping("/hello")
    public String hello() {
        return "hello world";
    }

    static AtomicInteger id = new AtomicInteger();

    @GetMapping("/invoke")
    public String invoke(@RequestParam(required = false, name = "service") String service,
                         @RequestParam(required = false, name = "parameter") String parameter,
                         @RequestParam(required = false, name = "content") String content,
                         @RequestParam(required = false, name = "host") String ipPort,
                         @RequestParam(required = false, name = "timeout") Integer timeout) {

        int port = 2045;
        String host = "127.0.0.1";

        // resolve ip port, connect to rpc server
        if (ipPort != null && ipPort.length() > 0) {
            // format: ip:port ?
            if (ipPort.contains(":")) {
                String[] hosts = ipPort.split(":");
                host = hosts[0];
                port = Integer.parseInt(hosts[1]);
            } else {
                host = ipPort;
            }
        }

        // resolve request parameters
        // format: key=value,another-key=another-value,...=...
        Map<String, String> requestParameters = new HashMap<>();
        if (parameter != null && parameter.length() > 0) {
            String[] items = parameter.split(",");
            for (String item : items) {
                if (item.contains("=")) {
                    String[] pair = item.split("=");
                    requestParameters.put(pair[0], pair[1]);
                }
            }
        }

        String defaultService = "com.alipay.core.UserService";
        service = (service == null || service.length() <= 0) ? defaultService : service;

        // inject required invoke interface name.
        requestParameters.put("interface", service);

        Protocol.Response response = BRpcClient.getClient(host, port)
                .request(id.incrementAndGet()
                        , requestParameters
                        , content == null ? "hello world" : content
                        , timeout == null ? 0 : timeout.intValue());

        return response.toString();
    }

}