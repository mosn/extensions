package com.alipay.plugin;

import org.apache.commons.logging.LogFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * @author yiji@apache.org
 */
@SpringBootApplication
public class XRpcServerBootStrap {

    static final Integer port = 7755;

    public static void main(String[] args) {

        // start netty server
        new Thread(() -> {
            try {
                XRpcController.main(args);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }).start();

        SpringApplication.run(XRpcServerBootStrap.class, args);

        printStartedInfo();
    }

    private static void printStartedInfo() {
        LogFactory.getLog(XRpcServerBootStrap.class).info("XRpc server initialized with port(s): " + port + " (tcp)");
    }

}