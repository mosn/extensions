## 插件编译:

编译命令执行都在go-plugin根项目中执行，详细请参考[插件基础篇](https://github.com/mosn/extensions/tree/master/go-plugin/doc/1.plugin-prepare.md), 编译插件语法：

```shell
make [module] plugin=[plugin-name]
# module: 
# 取值codec、filter和trans, 分别代表协议插件、拦截器插件和协议转换插件

# plugin-name:
# 代表插件的名称, 当module取值为filter或者trans时，允许以逗号分隔指定多个插件名称，脚手架同时编译多个插件
# 当module取值为filter是时，逗号分隔的插件名称代表拦截器执行的先后顺序
```

编译示例，比如编译bolt协议插件：

```shell
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

```shell
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

编译协议转换插件，以bolt转springcloud插件bolt2sp为例：

```shell
make trans plugin=bolt2sp
#  编译后在build/transcoders目录下输出so和配置
└── bolt2sp
    ├── egress_config.json
    ├── metadata.json
    ├── transcoder-bolt2sp.md5
    └── transcoder-bolt2sp.so

1 directory, 4 files
```

## 插件打包:

接下来介绍如何在本地打包插件代码，方便后续正式环境使用。和编译插件类似，以下打包命令执行都在go-plugin根项目中执行，打包插件语法：

```shell
make pkg-[module] plugin=[plugin-name]
# module: 
# 取值codec、filter和trans, 分别代表协议插件、拦截器插件和协议转换插件

# plugin-name:
# 代表插件的名称, 当module取值为filter或者trans时，允许以逗号分隔指定多个插件名称，脚手架同时打包多个插件
```

和编译插件区别，在模块前缀加了pkg前缀标识打包，用来将插件打包成.zip文件，用于在控制台上传。

比如打包bolt协议插件：

```shell
make pkg-codec plugin=bolt
# 编译后在build/target/codecs目录下输出bolt.zip
└── bolt.zip

0 directories, 1 file
```

打包拦截器插件, 以简单鉴权拦截器auth为例：

```shell
make pkg-filter plugin=auth
#  编译后在build/target/stream_filters目录下输出auth.zip
└── auth.zip

0 directories, 1 file
```

打包协议转换插件，以bolt转springcloud插件bolt2sp为例：

```shell
make pkg-trans plugin=bolt2sp
#  编译后在build/target/transcoders目录下输出bolt2sp.zip
└── bolt2sp.zip

0 directories, 1 file
```
