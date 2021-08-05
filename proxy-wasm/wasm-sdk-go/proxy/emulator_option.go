package proxy

type Option struct {
	pluginConfiguration ConfigMap
	vmConfiguration     ConfigMap
	newRootContext      func(uint32) RootContext
	newStreamContext    func(uint32, uint32) StreamContext
	newFilterContext    func(uint32, uint32) FilterContext
	newProtocolContext  func(uint32, uint32) ProtocolContext
}

func NewEmulatorOption() *Option {
	return &Option{}
}

func (o *Option) WithNewRootContext(f func(uint32) RootContext) *Option {
	o.newRootContext = f
	return o
}

func (o *Option) WithNewHttpContext(f func(uint32, uint32) FilterContext) *Option {
	o.newFilterContext = f
	return o
}

func (o *Option) WithNewStreamContext(f func(uint32, uint32) StreamContext) *Option {
	o.newStreamContext = f
	return o
}

func (o *Option) WithNewProtocolContext(f func(uint32, uint32) ProtocolContext) *Option {
	o.newProtocolContext = f
	return o
}

func (o *Option) WithPluginConfiguration(data ConfigMap) *Option {
	o.pluginConfiguration = data
	return o
}

func (o *Option) WithVMConfiguration(data ConfigMap) *Option {
	o.vmConfiguration = data
	return o
}
