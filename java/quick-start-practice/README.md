文档修订历史

| 版本号 | 作者 | 备注     | 修订日期      |
| ------ | ---- | -------- |-----------|
| 0.1    | [诣极](https://github.com/zonghaishang) | 初始版本 | 2022.5.20 |

## 1. 编译调试

### 1.1 编译mosn

有2种途径得到mosn可执行文件, 通过源码编译或者使用预编译好的压缩包文件。

**推荐：**
获取预编译商业版`mosn.zip`解压缩放入到`build/sidecar/binary`中。

解压缩后的mosn文件结构：

```shell
build
├── sidecar
   └── binary
       ├── mosn
       └── mosn-1.26.0-5ecd61124.md5
```

**说明**：商业版mosn.zip不允许放置开源，可以联系诣极或者邀明获取mosn.zip文件。

或者在脚手架根目录编译(依赖[商业版mosn源码](https://github.com/mosn/extensions/blob/master/go-plugin/doc/1.plugin-prepare.md#111-mosn%E6%BA%90%E7%A0%81) )：
```bash
# 1. 项目根目录执行mosn编译, 编译后在build/sidecar/binary目录自动输出可执行文件mosn
make ant
```

### 1.2 编译插件

以下编译命令执行都在go-plugin根项目中执行，编译插件语法：

```bash
make [codec|filter|trans] plugin=[plugin-name]
# 取值codec: 代表编译协议插件
# 取值filter: 代表编译拦截器插件
# 取值trans: 代表编译协议转换插件

# plugin-name:
# 代表插件的名称, 当编译类型取值为filter或者trans时，允许以逗号分隔指定多个插件名称，脚手架同时编译多个插件
# 当编译类型取值为filter时，逗号分隔的插件名称代表拦截器执行的先后顺序
```

编译示例，比如编译bolt协议插件：

```bash
make codec plugin=bolt
# 编译后在build/codecs目录下输出so和配置
└── bolt
    ├── codec-bolt.md5
    ├── codec-bolt.so
    ├── egress_bolt.json
    ├── ingress_bolt.json
    ├── metadata.json
    └── mosn_config.json

1 directory, 6 files
```

编译拦截器插件, 以简单鉴权拦截器auth为例：

```bash
make filter plugin=auth
#  编译后在build/stream_filters目录下输出so和配置
└── auth
    ├── egress_config.json
    ├── filter-auth.md5
    ├── filter-auth.so
    ├── metadata.json
    └── mosn_config.json

1 directory, 5 files
```

编译协议转换插件，以bolt转springcloud插件bolt2springcloud为例：

```bash
make trans plugin=bolt2springcloud
#  编译后在build/transcoders目录下输出so和配置
└── bolt2springcloud
    ├── egress_config.json
    ├── metadata.json
    ├── transcoder-bolt2springcloud.md5
    └── transcoder-bolt2springcloud.so

1 directory, 4 files
```

可以组合编译，让拦截器或者协议转换插件挂载到指定协议生效，比如：

```bash
# 1. 把拦截器插件auth挂载到bolt协议扩展中
make codec plugin=bolt filter=auth

# 2. 把协议转换插件bolt2springcloud挂载到bolt协议扩展中
make codec plugin=bolt trans=bolt2springcloud

# 3. 也允许同时挂载拦截器auth插件、协议转换插件bolt2springcloud到bolt协议扩展中
make codec plugin=bolt filter=auth trans=bolt2springcloud
```

注意：示例中auth插件要求请求头中包含`User`:`admin` 信息，否则请求会被拦截。

为了精简编译命令，比如make codec plugin=bolt filter=auth trans=bolt2springcloud可以同时完成auth插件、bolt2springcloud插件编译。

如果需要清除所有已经编译的插件，执行clean命令：

```bash
make clean
```

**交叉编译支持：**

- 支持显示指定目标交叉参数(arm64/amd64)

```
GOOS=linux GOARCH=amd64 make codec plugin=bolt
```

或者另外一种简写：

```shell
make codec.amd64 plugin=bolt
```

> apple m1 chip(arm64) 打包阶段自动触发编译amd64架构插件

在苹果mac m1 芯片一下命令，打包阶段自动重新编译linux amd64架构插件，方便公共一套linux sidecar容器。
```shell
make pkg-codec plugin=bolt
```

如果需要明确交叉编译插件，比如arm64插件，可以执行以下命令：

```shell
make pkg-codec.arm64 plugin=bolt
```

该命令等价于:

```shell
GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-codec plugin=bolt
```

- DOCKER_BUILD_OPTS 是传递给docker的编译参数，需要和传递给go `os/arch`架构一致。

### 1.3 编译调试

以下命令执行都在go-plugin根项目中执行，启动mosn调试命令：

```bash
make debug
# 启动后，终端输出以下提示信息：
# start mosn-container container success.
# run 'docker exec -it mosn-container /bin/bash' command enter mosn container.
```

注意：make debug 之后，mosn会监听2345端口，并且阻塞直到有远程调试行为, 接下来会介绍。make debug可以反复执行，重复执行会自动停止正在运行的mosn，nice~。

如果需要调试商业版mosn的源代码，可以把商业版mosn的源代码配置到脚手架， 如果要调试开源mosn，挂载开源mosn源码目录即可。操作步骤：

1. 在go-plugin项目中(`goland菜单`-> `Perferences...`->`Project Structure`-> `Add Content Root`)，选择商业版mosn代码根目录，如下图所示：

![image.png](./doc/images/add_mosn_src.png)

2. 在goland中新建一个远程调试配置， `Edit Configuratons...` -> `Go remote`，On Disconnect选择`Leave it runing`。

   ![远程调试](./doc/images/remote_debug.png)

3. 配置好远程调试后，在插件代码打好断点(在2.1小节介绍bolt调试)，点击Debug mosn按钮即可

![调试mosn](./doc/images/debug-mosn.png)

提示：为了提示本地排查效率，go-plugin根目录挂载了logs容器日志目录, 方便不进容器查看容器日志, nice~。

**说明：** 如果无需调试mosn和插件，可以执行启动命令：
```shell
make start
# 启动后，终端输出以下提示信息：
# start mosn-container container success.
# run 'docker exec -it mosn-container /bin/bash' command enter mosn container.
```

### 1.4 插件打包

接下来介绍如何在本地打包插件代码，方便后续正式环境使用。和编译插件类似，以下打包命令执行都在go-plugin根项目中执行，打包插件语法：

```bash
make pkg-[codec|filter|trans] plugin=[plugin-name]
# 取值codec: 代表打包协议插件
# 取值filter: 代表打包拦截器插件
# 取值trans: 代表打包协议转换插件

# plugin-name:
# 代表插件的名称, 当打包类型是filter或者trans时，允许以逗号分隔指定多个插件名称，脚手架同时打包多个插件
```

和编译插件区别，在模块前缀加了pkg前缀标识打包，用来将插件打包成.zip文件，用于在控制台上传。

比如打包bolt协议插件：

```bash
make pkg-codec plugin=bolt
# 编译后在build/target/codecs目录下输出bolt.zip
└── bolt.zip

0 directories, 1 file
```

打包拦截器插件, 以简单鉴权拦截器auth为例：

```bash
make pkg-filter plugin=auth
#  编译后在build/target/stream_filters目录下输出auth.zip
└── auth.zip

0 directories, 1 file
```

打包协议转换插件，以bolt转springcloud插件bolt2springcloud为例：

```bash
make pkg-trans plugin=bolt2springcloud
#  编译后在build/target/transcoders目录下输出bolt2springcloud.zip
└── bolt2springcloud.zip

0 directories, 1 file
```
