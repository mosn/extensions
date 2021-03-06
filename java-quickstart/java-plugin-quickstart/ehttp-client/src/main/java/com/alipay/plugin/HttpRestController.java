package com.alipay.plugin;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import javax.annotation.Resource;
import java.util.Arrays;

/**
 * @author yiji@apache.org
 */
@RestController()
public class HttpRestController {

    @Resource
    RestTemplate restTemplate;

    /**
     * example:
     * <p>
     * curl localhost:8008/hello
     * <p>
     * change request backend port(eg: request mosn port ?)
     * <p>
     * curl localhost:8008/hello?port=3045
     *
     * @param port
     * @return
     */
    // example : https://www.tutorialspoint.com/spring_boot/spring_boot_rest_template.htm
    @GetMapping("/hello")
    public String hello(@RequestParam(required = false) String port, @RequestParam(required = false) String user) {

        HttpHeaders headers = new HttpHeaders();

        // inject data id
        headers.set("X-SERVICE", "ehttp-provider");
        headers.setAccept(Arrays.asList(MediaType.APPLICATION_JSON));

        if (user != null && user.length() > 0) {
            headers.set("user", user);
        }

        String url = port != null && port.length() > 0
                ? "http://localhost:" + port + "/hello"
                : "http://localhost:3045/hello";

        return restTemplate.exchange(url
                        , HttpMethod.GET
                        , new HttpEntity<>(headers)
                        , String.class)
                .getBody();
    }

    @GetMapping("/userInfo")
    public String userInfo(@RequestParam(required = false) String port, @RequestParam(required = false) String user) {
        HttpHeaders headers = new HttpHeaders();

        // not inject header
        headers.setAccept(Arrays.asList(MediaType.APPLICATION_JSON));
        headers.setContentType(MediaType.APPLICATION_JSON);

        if (user != null && user.length() > 0) {
            headers.set("user", user);
        }

        String url = port != null && port.length() > 0
                ? "http://localhost:" + port + "/userInfo"
                : "http://localhost:3045/userInfo";


        String httpBody = "{\"target\": \"ehttp-provider\"}";
        HttpEntity<String> entity = new HttpEntity<String>(httpBody, headers);

        return restTemplate.exchange(url
                        , HttpMethod.POST
                        , entity
                        , String.class)
                .getBody();
    }
}