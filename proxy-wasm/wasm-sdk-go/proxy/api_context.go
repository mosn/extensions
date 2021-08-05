package proxy

const (
	CopyRawBytesFlag byte = 1           // 0000_0001
	HeartBeatFlag    byte = 1 << 5      // 0010_0000
	RpcRequestFlag   byte = 1 << 6      // 0100_0000
	RpcOnewayFlag    byte = 1<<6 | 1<<7 // 1100_0000
)

//export proxy_on_context_create
func proxyOnContextCreate(contextID uint32, rootContextID uint32) {
	if rootContextID == 0 {
		this.createRootContext(contextID)
	} else if this.newFilterContext != nil {
		this.createFilterContext(contextID, rootContextID)
	} else if this.newProtocolContext != nil {
		this.createProtocolContext(contextID, rootContextID)
	} else if this.newStreamContext != nil {
		this.createStreamContext(contextID, rootContextID)
	} else {
		panic("invalid context id on proxy_on_context_create")
	}
}

//export proxy_on_done
func proxyOnDone(contextID uint32) bool {
	if ctx, ok := this.filterStreams[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.OnFilterStreamDone()
		return true
	} else if _, ok := this.protocolStreams[contextID]; ok {
		this.setActiveContextID(contextID)
		//ctx.OnProtocolDone()
		return true
	} else if ctx, ok := this.streams[contextID]; ok {
		this.setActiveContextID(contextID)
		ctx.OnStreamDone()
		return true
	} else if ctx, ok := this.rootContexts[contextID]; ok {
		this.setActiveContextID(contextID)
		response := ctx.context.OnVMDone()
		return response
	} else {
		panic("invalid context on proxy_on_done")
	}
}

//export proxy_on_delete
func proxyOnDelete(contextID uint32) {
	if _, ok := this.filterStreams[contextID]; ok {
		delete(this.filterStreams, contextID)
	} else if _, ok := this.protocolStreams[contextID]; ok {
		delete(this.protocolStreams, contextID)
		// debug info
		//buf := bytes.Buffer{}
		//for id, _ := range this.protocolStreams {
		//	buf.WriteString(strconv.Itoa(int(id)))
		//	buf.WriteString(",")
		//}
		//log.Infof("delete contextId %d, remain: %s", contextID, buf.String())
	} else if _, ok := this.streams[contextID]; ok {
		delete(this.streams, contextID)
	} else if _, ok = this.rootContexts[contextID]; ok {
		delete(this.rootContexts, contextID)
	} else {
		panic("invalid context on proxy_on_delete")
	}
	delete(this.contextIDToRootID, contextID)
}
