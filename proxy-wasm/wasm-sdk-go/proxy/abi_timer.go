package proxy

//export proxy_on_tick
func proxyOnTick(rootContextID uint32) {
	ctx, ok := this.rootContexts[rootContextID]
	if !ok {
		panic("invalid root_context_id")
	}
	ctx.context.OnTick()
}
