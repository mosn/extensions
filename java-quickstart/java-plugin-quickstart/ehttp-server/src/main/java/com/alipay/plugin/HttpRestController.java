package com.alipay.plugin;

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yiji@apache.org
 */
@RestController()
public class HttpRestController {

    @GetMapping("/hello")
    public String hello() {
        return "hello world";
    }

    @PostMapping("/userInfo")
    @ResponseBody
    public Result userInfo(@RequestBody UserModel user) {

        Result result = new Result();
        result.userId = (user.userId == null || user.userId.length() <= 0) ? "yiji" : user.userId;
        if (user.parameters != null) {
            result.parameters = new HashMap<>();
            result.parameters.putAll(user.parameters);
        }

        return result;
    }

    public static class UserModel/* implements Serializable*/ {
        private String userId;

        private Map<String, String> parameters;

        public String getUserId() {
            return userId;
        }

        public void setUserId(String userId) {
            this.userId = userId;
        }

        public Map<String, String> getParameters() {
            return parameters;
        }

        public void setParameters(Map<String, String> parameters) {
            this.parameters = parameters;
        }
    }

    public static class Result/* implements Serializable*/ {
        private String userId;

        private String title = "developer";
        private String address = "hangzhou";

        private String responseIp;

        private Map<String, String> parameters;

        public String getUserId() {
            return userId;
        }

        public void setUserId(String userId) {
            this.userId = userId;
        }

        public Map<String, String> getParameters() {
            return parameters;
        }

        public void setParameters(Map<String, String> parameters) {
            this.parameters = parameters;
        }

        public String getTitle() {
            return title;
        }

        public void setTitle(String title) {
            this.title = title;
        }

        public String getAddress() {
            return address;
        }

        public void setAddress(String address) {
            this.address = address;
        }

        public String getResponseIp() {
            return responseIp;
        }

        public void setResponseIp(String responseIp) {
            this.responseIp = responseIp;
        }

    }
}