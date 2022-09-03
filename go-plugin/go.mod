module mosn.io/extensions/go-plugin

go 1.14

require (
	github.com/apache/dubbo-go-hessian2 v1.9.2 // dubbo
	github.com/dubbogo/gost v1.11.25 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.7.1
	github.com/valyala/fasthttp v1.31.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	mosn.io/api v0.0.0-20211217011300-b851d129be01
	mosn.io/pkg v0.0.0-20211217101631-d914102d1baf
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.4.3
	github.com/klauspost/compress => github.com/klauspost/compress v1.13.5

	github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	google.golang.org/grpc => google.golang.org/grpc v1.38.0
)
