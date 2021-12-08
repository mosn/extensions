module github.com/mosn/extensions/go-plugin

go 1.14

require (
	github.com/apache/dubbo-go-hessian2 v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.31.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	mosn.io/api v0.0.0-20211118092229-0f48ccc614b6
	mosn.io/pkg v0.0.0-20211019125153-96b01e984d62
)

replace (
	github.com/apache/dubbo-go-hessian2 => github.com/apache/dubbo-go-hessian2 v1.4.1-0.20200516085443-fa6429e4481d
	github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	mosn.io/api => github.com/YIDWang/api v0.0.0-20211119033852-6e898f4d8222
	mosn.io/pkg => mosn.io/pkg v0.0.0-20211208031317-7b020e8bcdd4
)
