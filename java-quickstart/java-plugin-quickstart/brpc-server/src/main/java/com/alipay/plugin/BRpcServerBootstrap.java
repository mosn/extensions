package com.alipay.plugin;

import org.apache.commons.logging.LogFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * @author yiji@apache.org
 */
@SpringBootApplication
public class BRpcServerBootstrap {

    static final Integer port = 7766;

    public static void main(String[] args) {

        // start netty server
        new Thread(() -> {
            try {
                BRpcController.main(args);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }).start();

        SpringApplication.run(BRpcServerBootstrap.class, args);

        printStartedInfo();
    }

    private static void printStartedInfo() {
        LogFactory.getLog(BRpcServerBootstrap.class).info("BRpc server initialized with port(s): " + port + " (tcp)");
    }
}