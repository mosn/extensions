文档修订历史

| 版本号 | 作者 | 备注                                                                             | 修订日期  |
| ------ | ---- |--------------------------------------------------------------------------------| --------- |
| 0.1    | [诣极](https://github.com/zonghaishang) | 初始版本                                                                           | 2022.1.12 |
| 0.2 | [诣极](https://github.com/zonghaishang) | 新增[扩展点api介绍](https://github.com/mosn/extensions/tree/master/go-plugin/doc/api) | 2022.4.6 |

## [1. 插件基础篇](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#1-%E6%8F%92%E4%BB%B6%E5%9F%BA%E7%A1%80%E7%AF%87)

### [1.1 环境准备](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#11-%E7%8E%AF%E5%A2%83%E5%87%86%E5%A4%87)

#### [1.1.1 mosn源码](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#111-mosn%E6%BA%90%E7%A0%81)

#### [1.1.2 插件源码](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#112-%E6%8F%92%E4%BB%B6%E6%BA%90%E7%A0%81)

#### [1.1.3 插件介绍](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#113-%E6%8F%92%E4%BB%B6%E4%BB%8B%E7%BB%8D)

### [1.2 编译调试](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#12-%E7%BC%96%E8%AF%91%E8%B0%83%E8%AF%95)

#### [1.2.1 编译mosn](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#121-%E7%BC%96%E8%AF%91mosn)

#### [1.2.2 编译插件](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#122-%E7%BC%96%E8%AF%91%E6%8F%92%E4%BB%B6)

#### [1.2.3 编译调试](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#123-%E7%BC%96%E8%AF%91%E8%B0%83%E8%AF%95)

#### [1.2.4 插件打包](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#124-%E6%8F%92%E4%BB%B6%E6%89%93%E5%8C%85)

## [2. mesh功能扩展篇](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md)

### [2.1 动手实现bolt协议插件化](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#21-%E5%8A%A8%E6%89%8B%E5%AE%9E%E7%8E%B0bolt%E5%8D%8F%E8%AE%AE%E6%8F%92%E4%BB%B6%E5%8C%96)

#### [2.1.1 编解码实现](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#211-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AE%9E%E7%8E%B0)

#### [2.1.2 编解码对象](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#212-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AF%B9%E8%B1%A1)

#### [2.1.3 心跳处理](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#213-%E5%BF%83%E8%B7%B3%E5%A4%84%E7%90%86) 

#### [2.1.4 请求劫持](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#214-%E8%AF%B7%E6%B1%82%E5%8A%AB%E6%8C%81) 

#### [2.1.5 协议Codec](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#215-%E5%8D%8F%E8%AE%AEcodec) 

#### [2.1.6 启动bolt应用服务](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.1bolt.md#216-%E5%90%AF%E5%8A%A8bolt%E5%BA%94%E7%94%A8%E6%9C%8D%E5%8A%A1) 

### [2.2 动手实现标准dubbo协议扩展](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md)

#### [2.2.1 编解码实现](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md#221-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AE%9E%E7%8E%B0)

#### [2.2.2 编解码对象](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md#222-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AF%B9%E8%B1%A1)

#### [2.2.3 心跳处理](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md#223-%E5%BF%83%E8%B7%B3%E5%A4%84%E7%90%86)

#### [2.2.4 请求劫持](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md#224-%E8%AF%B7%E6%B1%82%E5%8A%AB%E6%8C%81)

#### [2.2.5 协议Codec](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md#225-%E5%8D%8F%E8%AE%AEcodec)

#### [2.2.6 启动dubbo应用服务](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.2dubbo.md#226-%E5%90%AF%E5%8A%A8dubbo%E5%BA%94%E7%94%A8%E6%9C%8D%E5%8A%A1)

### [2.3 传统xml协议标准接入实战](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md) 

#### [2.3.1 编解码实现](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md#231-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AE%9E%E7%8E%B0) 

#### [2.3.2 编解码对象](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md#232-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AF%B9%E8%B1%A1) 

#### [2.3.3 心跳处理](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md#233-%E5%BF%83%E8%B7%B3%E5%A4%84%E7%90%86) 

#### [2.3.4 请求劫持](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md#234-%E8%AF%B7%E6%B1%82%E5%8A%AB%E6%8C%81) 

#### [2.3.5 协议Codec](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md#235-%E5%8D%8F%E8%AE%AEcodec) 

#### [2.3.6 启动xr应用服务](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.3xml.md#236-%E5%90%AF%E5%8A%A8xr%E5%BA%94%E7%94%A8%E6%9C%8D%E5%8A%A1) 

### [2.4 轻松实现http协议扩展](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md)

#### [2.4.1 编解码实现](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md#241-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AE%9E%E7%8E%B0)

#### [2.4.2 编解码对象](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md#242-%E7%BC%96%E8%A7%A3%E7%A0%81%E5%AF%B9%E8%B1%A1)

#### [2.4.3 协议Codec](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md#243-%E5%8D%8F%E8%AE%AEcodec)

#### [2.4.4 启动springcloud应用服务](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md#244-%E5%90%AF%E5%8A%A8springcloud%E5%BA%94%E7%94%A8%E6%9C%8D%E5%8A%A1)

#### [2.4.5 再谈http协议扩展](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md#245-%E5%86%8D%E8%B0%88http%E5%8D%8F%E8%AE%AE%E6%89%A9%E5%B1%95)

#### [2.4.6 获取http服务标识](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.4springcloud.md#246-%E8%8E%B7%E5%8F%96http%E6%9C%8D%E5%8A%A1%E6%A0%87%E8%AF%86)

### [2.5 拦截器](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.5filter.md) 

#### [2.5.1 动手实现http简易鉴权拦截器](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.5filter.md#251-%E5%8A%A8%E6%89%8B%E5%AE%9E%E7%8E%B0http%E7%AE%80%E6%98%93%E9%89%B4%E6%9D%83%E6%8B%A6%E6%88%AA%E5%99%A8) 

#### [2.5.2 启动springcloud应用服务](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.5filter.md#252-%E5%90%AF%E5%8A%A8springcloud%E5%BA%94%E7%94%A8%E6%9C%8D%E5%8A%A1)  

### [2.6 协议转换插件](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.6transcoder.md)

#### [2.6.1 标准dubbo和spring cloud协议互转实践](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.6.1dubbo2springcloud.md)

#### [2.6.2 标准bolt和spring cloud协议互转实践](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.6.2bolt2springcloud.md) 

#### [2.6.3 传统xml和spring cloud协议互转实践](https://github.com/mosn/extensions/blob/master/go-plugin/doc/2.6.3xml2springcloud.md) 

2.7 mesh治理能力

2.7.1 服务限流能力

2.7.1.1 标准dubbo接入服务限流能力

2.7.1.2 传统xml协议接入服务限流能力

2.7.1.3 标准springcloud接入服务限流能力

2.7.2 服务熔断能力

2.7.2.1 标准dubbo接入服务熔断能力

2.7.2.2 传统xml协议接入服务熔断能力

2.7.3 服务降级能力

2.7.3.1 标准dubbo接入服务降级能力

2.7.3.2 传统xml协议接入服务降级能力

2.7.4 服务故障注入能力

2.7.4.1 标准dubbo接入服务降级能力

2.7.4.2 传统xml协议接入服务降级能力

## [3.mesh扩展使用篇](https://github.com/mosn/extensions/blob/master/go-plugin/doc/3.1pub-pod.md) 

### [3.1 发布容器服务](https://github.com/mosn/extensions/blob/master/go-plugin/doc/3.1pub-pod.md#31-%E5%8F%91%E5%B8%83%E5%AE%B9%E5%99%A8%E6%9C%8D%E5%8A%A1)

#### [3.1.1 手动发布容器应用](https://github.com/mosn/extensions/blob/master/go-plugin/doc/3.1pub-pod.md#311-%E5%AE%B9%E5%99%A8%E4%B8%AD%E5%8F%91%E5%B8%83%E5%BA%94%E7%94%A8)

4. 开放扩展原理篇

4.1 标准插件扩展原理

4.2 理解传统xml接入原理

4.2.1 为什么传统xml协议需要每个tcp连接对应一个协议实例

4.3 理解http协议扩展原理

4.4 拦截器插件原理

4.5 协议转换插件原理

