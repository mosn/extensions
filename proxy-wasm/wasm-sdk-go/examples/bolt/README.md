

# 创建插件工程

### 1.1 前置准备

- 安装[go](https://golang.org/doc/install)
- 安装[tinygo](https://tinygo.org/getting-started/linux/) 

提示： 如果已有go不需要重复安装，tinygo用于编译成wasm插件

### 1.2 创建项目工程

在$GOPATH/src目录中创建工程, 假设项目工程名为`plugin-repo`, 用来包含多个插件：

```shell
# 1. 查看GOPATH路径
go env | grep GOPATH

# 2. 在GOPATH/src目录中创建
mkdir plugin-repo
cd plugin-repo

# 3. 执行项目初始化
go mod init

# 4. 创建协议插件目录名称，假设叫做bolt
mkdir -p bolt/main
```

执行完成后，目录结构如下：

```go
plugin-repo				// 插件仓库根目录
├── go.mod				// 项目依赖管理
└── bolt					// 插件名称，开发者扩展代码放到这里
    └── main			// 注册插件逻辑，开发者编写注册插件逻辑
		└── build			// 插件编译后，自动生成
```

因为在开始编写插件时，需要依赖wasm sdk，需要在插件根目录，执行以下命令，拉取依赖：

```shell
 go get github.com/zonghaishang/wasm-sdk-go
 go mod vendor
```

提示：完整实例程序已经包含在github仓库, 请参考[plugin-repo](https://github.com/zonghaishang/plugin-repo)。

### 1.3 编写插件扩展

在开始编写插件前，我们先展示编写完成后的目录结构：

```go
plugin-repo				// 插件仓库根目录
├── go.mod				// 项目依赖管理
├── Makefile			// 编译插件成wasm文件
└── bolt
    ├── protocol.go
    ├── command.go
    ├── codec.go
    ├── api.go
    ├── main
    │   ├── main.go
    │   └── main_test.go
    ├── build			// 插件编译后，自动生成
        └── bolt-go.wasm
```

在协议扩展场景，我们主要提供编解码(codec)、编解码对象(command)、协议层支持心跳/响应(protocol)和注册协议(main)。



#### 1.3.1 编解码实现

在处理请求和响应流程中，开发者需要实现`Codec`接口, 方法处理逻辑如下：

- Decode：需要开发者将data中的字节数据解码成请求或者响应 
- Encode：需要开发者将请求或者响应编码成字节buffer

```go
type Codec interface {
	Decode(ctx context.Context, data Buffer) (Command, error)
	Encode(ctx context.Context, cmd Command) (Buffer, error)
}
```

注意：在Decode流程中，完成解码需要调用`data.Drain(frameLen)`, `frameLen`代表完整请求或者报文总长度

开发者在编写编解码时，建议采用协议名+Codec命名，比如bolt编解码，命名为`boltCodec`。

目前提供了示例编解码实现，请参考[boltCodec](https://github.com/zonghaishang/plugin-repo/blob/master/bolt/codec.go) .

#### 1.3.2 编解码对象

编解码主要在二进制字节流和请求/响应对象互转，开发者在定义请求/响应对象，应该遵守command接口。目前command主要分2类，请求和响应。

- 请求对象： 主要包括请求(request-response)、请求(oneway)、心跳类型
- 响应对象： 主要包括响应请求结果对象

请求对象除了表达request-response模型、oneway和心跳，也会承载超时等熟悉，与之对应响应会承载响应状态码。

目前请求和响应的接口契约如下：

```go
type Request interface {
	Command
	// IsOneWay Check that the request does not care about the response
	IsOneWay() bool
	GetTimeout() uint32 // request timeout
}

type Response interface {
	Command
	GetStatus() uint32 // response status
}
```

不管请求还是响应，除了识别command类型，还承担请求头部和请求体2部分，头部是普通的key-value结构，data部分应该是协议的content部分，而不是完整报文内容。

目前command的接口定义如下：

```go
// Command base request or response command
type Command interface {
	// Header get the data exchange header, maybe return nil.
	GetHeader() Header
	// GetData return the full message buffer, the protocol header is not included
	GetData() Buffer
	// SetData update the full message buffer, the protocol header is not included
	SetData(data Buffer)
	// IsHeartbeat check if the request is a heartbeat request
	IsHeartbeat() bool
	// CommandId get command id
	CommandId() uint64
	// SetCommandId update command id
	// In upstream, because of connection multiplexing,
	// the id of downstream needs to be replaced with id of upstream
	// blog: https://mosn.io/blog/posts/multi-protocol-deep-dive/#%E5%8D%8F%E8%AE%AE%E6%89%A9%E5%B1%95%E6%A1%86%E6%9E%B6
	SetCommandId(id uint64)
}
```

目前提供了示例编解码对象实现，请参考[command](https://github.com/zonghaishang/plugin-repo/blob/master/bolt/command.go).

#### 1.3.3 协议层

因为心跳需要协议层理解，如果开发者扩展的协议支持心跳能力，应当提供扩展`KeepAlive`实现：

- KeepAlive: 根据请求id生成一个心跳请求command
- ReplyKeepAlive: 根据收到的请求，返回一个心跳响应command

```go
type KeepAlive interface {
	KeepAlive(requestId uint64) Request
	ReplyKeepAlive(request Request) Response
}
```

注意： 如果扩展协议不支持心跳或者不需要心跳，协议层KeepAlive方法返回nil即可

在service mesh场景中，因为增加了一跳，mesh在转发过程中可能被控制面拦截，比如限流熔断，需要协议层构造并返回响应，因此开发者需要提供`Hijacker`接口实现：

- Hijack: 根据请求和拦截状态码，返回一个响应command

```go
type Hijacker interface {
	// Hijack allows sidecar to hijack requests
	Hijack(request Request, code uint32) Response
}
```

目前协议层接口采用组合方式，主要讲编解码独立拆分出去， protocol接口定义：

```go
type Protocol interface {
	Name() string
	Codec() Codec
	KeepAlive
	Hijacker
	Options
}
```

接口中方法描述：

- Name： 返回协议名称
- Codec：返回协议编解码对象，对应`1.3.1`小节
- KeepAlive：协议心跳实现
- Hijacker： 处理控制面拦截逻辑
- Options： 协议层配置选项开发，一般协议组合默认配置`proxy.DefaultOptions`

目前提供了示例协议实现，请参考[protocol](https://github.com/zonghaishang/plugin-repo/blob/master/bolt/protocol.go).

#### 1.3.4 注册协议

在完成协议扩展后，需要将我们编写的插件进行注册，在wasm扩展中，我们一切是以Context为核心来转的，比如host侧触发解码，在沙箱内会调用开发者protocol context的回调来解码。

因此注册协议我们需要提供一个`ProtocolContext`接口实现，和protocol接口极其类似:

```go
// L7 layer extension
type ProtocolContext interface {
	Name() string         // protocol name
	Codec() Codec         // frame encode & decode
	KeepAlive() KeepAlive // protocol keep alive
	Hijacker() Hijacker   // protocol hijacker
	Options() Options     // protocol options
}
```

以bolt协议插件为例，我们提供`boltProtocolContext`实现: 

```go
// 1. 提供bolt插件protocolContext实现
type boltProtocolContext struct {
	proxy.DefaultRootContext 			// notify on plugin start.
	proxy.DefaultProtocolContext 	// 继承默认协议实现，比如使用默认Options()
	bolt      proxy.Protocol			// 插件真实协议实现
	contextID uint32
}

// 2. 创建bolt单实例协议实例
var boltProtocol = bolt.NewBoltProtocol()

func boltContext(rootContextID, contextID uint32) proxy.ProtocolContext {
	return &boltProtocolContext{
		bolt:      boltProtocol,
		contextID: contextID,
	}
}

// 3. 注册boltContext协议钩子
func main() {
	proxy.SetNewProtocolContext(boltContext)
}

// 4. 如果协议不支持心跳，这里允许返回nil
func (proto *boltProtocolContext) KeepAlive() proxy.KeepAlive {
	return proto.bolt
}

// 5. 如果需要获取插件参数，可以override对应方法
func (proto *boltProtocolContext) OnPluginStart(conf proxy.ConfigMap) bool {
	proxy.Log.Infof("proxy_on_plugin_start from Go!")
	return true
}
```

目前提供了示例协议注册实现，请参考[main](https://github.com/zonghaishang/plugin-repo/blob/master/bolt/main/main.go) .

#### 1.3.5 调试&打包

开发者在编写完插件后，允许在本地idea直接开始调试测试，并且不依赖mosn启动。目前推荐在协议开发完后，提供main_test.go实现，在里面写集成测试。

目前wasm sdk提供了模拟器实现(`Emulator`), 可以模拟完整的mosn处理流程，并且可以回调开发者插件对应生命周期方法。基本用法：

```go
	// 1. 注册对应context和配置，boltContext在同一个main包下已经实现
	opt := proxy.NewEmulatorOption().
		WithNewProtocolContext(boltContext).
		WithNewRootContext(rootContext).
		WithVMConfiguration(vmConfig)

	// 2. 创建一个sidecar模拟器
	host := proxy.NewHostEmulator(opt)
	// release lock and reset emulator state
	defer host.Done()

	// 3. 调用host对应实现，比如启动沙箱
	host.StartVM()

	// 4. 调用启动插件
	host.StartPlugin()

	// 5. 模拟新请求到来，创建插件上下文
	ctxId := host.NewProtocolContext()

	// 6. 模拟host接收客户端请求，并解码
	cmd, err := host.Decode(...)

	// 7. 模拟host转发请求，并编码
	upstreamBuf, err := host.Encode(...)

	// 8. 模拟host处理完请求
	host.CompleteProtocolContext(ctxId)
```

如果要在`GoLand`中直接调试集成测试， 需要执行以下操作：

- `GoLand`->`Preferences...`->`Go`->`Build Tags & Vendoring`->`Custom tags`填写`proxytest`
- 调试窗口`Edit Configurations...`->勾选`Use all custom build tags`

目前提供了示例集成测试，请参考[main test](https://github.com/zonghaishang/plugin-repo/blob/master/bolt/main/main_test.go) .

目前打包插件，可以本地开发环境编译打包，也支持镜像方式编译插件, 目前通用[makefile](https://github.com/zonghaishang/plugin-repo/blob/master/Makefile)已经提供，可以copy到插件项目根目录中使用。

基于makefile，2种打包命令分别如下(编译成功会在插件中创建build文件夹，并且输出bolt-go.wasm)：

```
# 1. 本地编译，bolt替换成开发者插件名
make name=bolt

# 2. 基于镜像编译
make build-image name=name=bolt
```

目前提供了示例打包文件，请参考[bolt-go](https://github.com/zonghaishang/plugin-repo/tree/master/bolt/build) .

### 1.4 启动mosn

目前提供了一份用于wasm启动的配置文件[mosn_rpc_config_wasm.json](https://github.com/mosn/mosn/blob/master/configs/mosn_rpc_config_wasm.json),  可以使用以下命令启动mosn:

```shell
./mosnd start -c /path/to/mosn_rpc_config_wasm.json
```

提示：

- 目前提供的配置，会开启`2045`和`2046`端口，`2045`接收客户端请求，通过`2046`转发给服务端

- mosn_rpc_config_wasm中已经配置了bolt-go.wasm, 在项目根目录`etc/wasm/`目录中
- 如果是自定义协议插件，配置mosn_rpc_config_wasm.json中有几点需要修改
  - `vm_config.path`指向的wasm路径
  - `wasm_global_plugins.plugin_name`和`codecs.config.from_wasm_plugin`要相同
  - `codecs.config.from_wasm_plugin`和`extend_config.sub_protocol`要相同(一般协议有2个listener都要改)

其中， mosnd 可执行文件可以通过编译mosn获取, 执行以下命令:

```shell
# 下载mosn代码到本地GOPATH, 可以通过本地shell执行：go env | grep GOPATH 查看
# step 1: 
mkdir -p $GOPATH/src/mosn.io
cd $GOPATH/src/mosn.io

# step 2: 
# 因为当前merge request正在推进合并中，在fork开发分支上编译
git clone https://github.com/mosn/mosn.git

# step 3: 
# 本地编译
sudo make build-local

# 编译成功后，会在项目根目录下
build/bundles/v0.21.0/binary/mosnd
```

如果是研发同学，可以根据 `step 2` 拉取代码，直接通过`GoLand`右键项目根目录Debug(这样就不用手动去编译以及不需要命令行启动mosn了), 在`Edit Configurations...` 调试配置页签中修改包路径和程序入口参数：

```shell
Package path: mosn.io/mosn/cmd/mosn/main
Program arguments: start -c /path/to/mosn_rpc_config_wasm.json
```

提示：

- `/path/to` 需要替换成mosn根目录中到mosn_rpc_config_wasm.json文件的完整路径

### 1.5 启动应用服务

开发完成后，可以先启动mosn，然后启动应用的服务端和客户端，以sofaboot应用为例展示。

目前sofaboot应用测试程序已经托管到github上，可以通过以下命令获取:

```shell
git clone https://github.com/sofastack-guides/sofastack-mesh-demo.git 
# checkout到wasm_benchmark分支
git checkout wasm_benchmark

cd sofastack-mesh-demo/sofa-samples-springboot2
# 本地打包sofaboot应用程序
mvn clean package
# 打包成功后，会在sofa-echo-server和sofa-echo-client下生成target目录，
# 其中分别包含服务端和客户端可执行程序，文件名分别为：
# sofa-echo-server-web-1.0-SNAPSHOT-executable.jar
# sofa-echo-client-web-1.0-SNAPSHOT-executable.jar
```

启动sofaboot服务端程序：

```shell
java -DMOSN_ENABLE=true -Drpc_tr_port=12199 -Dspring.profiles.active=dev -Drpc_register_registry_ignore=true -jar sofa-echo-server-web-1.0-SNAPSHOT-executable.jar
```

然后启动sofaboot客户端程序：

```shell
java  -DMOSN_ENABLE=true -Drpc_tr_port=12198 -Dspring.profiles.active=dev -Drpc_register_registry_ignore=true -jar sofa-echo-client-web-1.0-SNAPSHOT-executable.jar
```

当客户端启动成功后，会在终端输出以下信息(每隔1秒发起一次wasm请求)：

```shell
>>>>>>>> [57,21,7ms]2021-03-16 20:57:05 echo result: Hello world!
>>>>>>>> [57,22,5ms]2021-03-16 20:57:06 echo result: Hello world!
>>>>>>>> [57,23,7ms]2021-03-16 20:57:07 echo result: Hello world!
>>>>>>>> [57,24,7ms]2021-03-16 20:57:08 echo result: Hello world!
>>>>>>>> [57,25,8ms]2021-03-16 20:57:09 echo result: Hello world!
>>>>>>>> [57,26,7ms]2021-03-16 20:57:10 echo result: Hello world!
>>>>>>>> [57,27,5ms]2021-03-16 20:57:11 echo result: Hello world!
>>>>>>>> [57,28,7ms]2021-03-16 20:57:12 echo result: Hello world!
```

当前扩展特性已经合并进开源社区，感兴趣同学可以查看实现原理：

- [wasm protocol #1579](https://github.com/mosn/mosn/pull/1597)
- [mosn api #31](https://github.com/mosn/api/pull/31)
- [wasm sdk-go](
