module mosn.io/extensions/go-plugin

go 1.18

require (
	github.com/SkyAPM/go2sky v0.5.0
	github.com/apache/dubbo-go-hessian2 v1.9.2 // dubbo
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.8.1
	github.com/valyala/fasthttp v1.40.0
	google.golang.org/grpc v1.53.0
	mosn.io/api v1.3.0
	mosn.io/pkg v1.3.1-0.20230104081003-9c747062328a
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dubbogo/gost v1.11.25 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/klauspost/compress => github.com/klauspost/compress v1.13.5
	github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
)
