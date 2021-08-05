package proxy

type (
	HttpCalloutCallBack = func(headers Header, body Buffer)

	rootContextState struct {
		context       RootContext
		httpCallbacks map[uint32]*struct {
			callback        HttpCalloutCallBack
			callerContextID uint32
		}
	}
)

type state struct {
	newRootContext     func(contextID uint32) RootContext
	rootContexts       map[uint32]*rootContextState
	newFilterContext   func(rootContextID, contextID uint32) FilterContext
	filterStreams      map[uint32]FilterContext
	newProtocolContext func(rootContextID, contextID uint32) ProtocolContext
	protocolStreams    map[uint32]ProtocolContext
	newStreamContext   func(rootContextID, contextID uint32) StreamContext
	streams            map[uint32]StreamContext

	// protocol context

	contextIDToRootID map[uint32]uint32
	activeContextID   uint32
}

var this = newState()

func newState() *state {
	return &state{
		rootContexts:      make(map[uint32]*rootContextState),
		filterStreams:     make(map[uint32]FilterContext),
		protocolStreams:   make(map[uint32]ProtocolContext),
		streams:           make(map[uint32]StreamContext),
		contextIDToRootID: make(map[uint32]uint32),
	}
}

func SetNewRootContext(f func(contextID uint32) RootContext) {
	this.newRootContext = f
}

func SetNewFilterContext(f func(rootContextID, contextID uint32) FilterContext) {
	this.newFilterContext = f
}

func SetNewStreamContext(f func(rootContextID, contextID uint32) StreamContext) {
	this.newStreamContext = f
}

func SetNewProtocolContext(f func(rootContextID, contextID uint32) ProtocolContext) {
	this.newProtocolContext = f
}

//go:inline
func (s *state) createRootContext(contextID uint32) {
	var ctx RootContext
	if s.newRootContext == nil {
		ctx = &DefaultRootContext{}
	} else {
		ctx = s.newRootContext(contextID)
	}

	s.rootContexts[contextID] = &rootContextState{
		context: ctx,
		httpCallbacks: map[uint32]*struct {
			callback        HttpCalloutCallBack
			callerContextID uint32
		}{},
	}
}

func (s *state) createFilterContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.filterStreams[contextID]; ok {
		panic("filter context id duplicated")
	}

	ctx := s.newFilterContext(rootContextID, contextID)
	s.contextIDToRootID[contextID] = rootContextID
	s.filterStreams[contextID] = ctx
}

func (s *state) createProtocolContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.protocolStreams[contextID]; ok {
		panic("protocol context id duplicated")
	}

	ctx := s.newProtocolContext(rootContextID, contextID)
	s.contextIDToRootID[contextID] = rootContextID
	s.protocolStreams[contextID] = ctx

}

func (s *state) createStreamContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.streams[contextID]; ok {
		panic(" stream context id duplicated")
	}

	ctx := s.newStreamContext(rootContextID, contextID)
	s.contextIDToRootID[contextID] = rootContextID
	s.streams[contextID] = ctx
}

func (s *state) registerHttpCallOut(calloutID uint32, callback HttpCalloutCallBack) {
	r := s.rootContexts[s.contextIDToRootID[s.activeContextID]]
	r.httpCallbacks[calloutID] = &struct {
		callback        HttpCalloutCallBack
		callerContextID uint32
	}{callback: callback, callerContextID: s.activeContextID}
}

//go:inline
func (s *state) setActiveContextID(contextID uint32) {
	s.activeContextID = contextID
}

func VMStateReset() {
	this = newState()
}

func VMStateGetActiveContextID() uint32 {
	return this.activeContextID
}
