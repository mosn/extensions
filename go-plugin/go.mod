module github.com/mosn/extensions/go-plugin

go 1.14

require (
	github.com/valyala/fasthttp v1.31.0
	mosn.io/api v0.0.0-20211217011300-b851d129be01
	mosn.io/pkg v0.0.0-20211217101631-d914102d1baf
)

replace github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // mosn dependency
