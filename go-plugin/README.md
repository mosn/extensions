文档修订历史

| 版本号 | 作者 | 备注     | 修订日期  |
| ------ | ---- | -------- | --------- |
| 0.1    | 诣极 | 初始版本 | 2022.1.12 |

# 1. [插件基础篇](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#1-%E6%8F%92%E4%BB%B6%E5%9F%BA%E7%A1%80%E7%AF%87)

## 1.1 [环境准备](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#11-%E7%8E%AF%E5%A2%83%E5%87%86%E5%A4%87)
### 1.1.1 [mosn源码](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#111-mosn%E6%BA%90%E7%A0%81)
### 1.1.2 [插件源码](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#112-%E6%8F%92%E4%BB%B6%E6%BA%90%E7%A0%81)

### 1.1.3 [插件介绍](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#113-%E6%8F%92%E4%BB%B6%E4%BB%8B%E7%BB%8D)

## 1.2 [编译调试](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#12-%E7%BC%96%E8%AF%91%E8%B0%83%E8%AF%95)

### 1.2.1 [编译mosn](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#121-%E7%BC%96%E8%AF%91mosn)

### 1.2.2 [编译插件](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#122-%E7%BC%96%E8%AF%91%E6%8F%92%E4%BB%B6)

### 1.2.3 [编译调试](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#123-%E7%BC%96%E8%AF%91%E8%B0%83%E8%AF%95)
### [1.2.4 插件打包](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#124-%E6%8F%92%E4%BB%B6%E6%89%93%E5%8C%85)

# 2. mesh功能扩展篇

2.1 动手实现bolt协议插件化
2.1.1 编解码实现
2.1.2 编解码对象
2.1.3 心跳处理
2.1.4 请求劫持
2.1.5 协议Codec
2.1.6 启动bolt应用服务

2.2 动手实现标准dubbo协议扩展
2.3 传统xml协议标准接入实战
2.4 轻松实现http协议扩展
2.5 拦截器
2.5.1 动手实现http简易鉴权拦截器

2.6 协议转换插件
2.6.1 标准dubbo和spring cloud协议互转实践
2.6.2 标准bolt和spring cloud协议互转实践
2.6.3 传统xml和spring cloud协议互转实践

2.7 mesh治理能力
2.7.1 服务限流能力
2.7.1.1 标准dubbo接入服务限流能力
2.7.1.2 传统xml协议接入服务限流能力

2.7.2 服务熔断能力
2.7.2.1 标准dubbo接入服务熔断能力
2.7.2.2 传统xml协议接入服务熔断能力

2.7.3 服务降级能力
2.7.3.1 标准dubbo接入服务降级能力
2.7.3.2 传统xml协议接入服务降级能力

2.7.4 服务故障注入能力
2.7.4.1 标准dubbo接入服务降级能力
2.7.4.2 传统xml协议接入服务降级能力

3. mesh扩展治理篇

3.1 如何在生产中使用插件
3.1.1 插件维护
3.1.1.1 功能插件维护
3.1.1.1.1 协议转换插件
3.1.1.1.2 拦截器插件
3.1.1.2 协议插件维护

3.1.2 sidecar注入规则
3.1.3 手动服务注册发现

3.2 如何在生产中激活治理能力

3.2.1 激活服务限流能力
3.2.1 激活服务熔断能力
3.2.3 激活服务降级能力
3.2.4 激活服务故障注入能力

4. 开放扩展原理篇

4.1 标准插件扩展原理
4.2 理解传统xml接入原理
4.2.1 为什么传统xml协议需要每个tcp连接对应一个协议实例
4.3 理解http协议扩展原理
4.4 拦截器插件原理
4.5 协议转换插件原理

