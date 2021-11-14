module github.com/mosn/wasm-sdk/go-plugin

go 1.14

require (
	mosn.io/api v0.0.0-20210714065837-5b4c2d66e70c
	mosn.io/pkg v0.0.0-20210823090748-f639c3a0eb36
)

//replace mosn.io/mosn v0.25.0 => github.com/zonghaishang/mosn v0.17.1-0.20211111054142-76358bb1e33d

replace mosn.io/api v0.0.0-20210714065837-5b4c2d66e70c => github.com/zonghaishang/api v0.0.0-20211111063821-9f6ab4c6e576
