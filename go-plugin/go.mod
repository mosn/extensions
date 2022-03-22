module mosn.io/extensions/go-plugin

go 1.14

require (
	github.com/beevik/etree v1.1.0
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.31.0
	github.com/valyala/fastjson v1.6.3
	mosn.io/api v0.0.0-20211217011300-b851d129be01
	mosn.io/pkg v0.0.0-20211217101631-d914102d1baf
)

replace (
	github.com/klauspost/compress => github.com/klauspost/compress v1.13.5
	github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // mosn dependency
)
