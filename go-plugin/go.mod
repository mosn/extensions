module github.com/mosn/extensions/go-plugin

go 1.14

require (
	github.com/valyala/fasthttp v1.31.0
	mosn.io/api v0.0.0-20211118092229-0f48ccc614b6
	mosn.io/pkg v0.0.0-20211019125153-96b01e984d62
)

replace mosn.io/pkg v0.0.0-20211019125153-96b01e984d62 => github.com/Tanc010/pkg v0.0.0-20211207033231-f2e906300a49
