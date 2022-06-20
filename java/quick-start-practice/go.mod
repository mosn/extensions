module github.com/zonghaishang/quick-start-practice

go 1.17

require (
	mosn.io/api v0.0.0-20211217011300-b851d129be01
	mosn.io/pkg v0.0.0-20211217101631-d914102d1baf
	github.com/valyala/fasthttp v1.31.0
)

replace (
	github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	github.com/klauspost/compress => github.com/klauspost/compress v1.13.5 // fast http
)
