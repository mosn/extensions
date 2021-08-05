// Copyright 2020 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/examples/bolt"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy"
)

func main() {
	proxy.SetNewRootContext(rootContext)
	proxy.SetNewProtocolContext(boltContext)
}

func rootContext(rootContextID uint32) proxy.RootContext {
	return &boltProtocolContext{
		bolt:      boltProtocol,
		contextID: rootContextID,
	}
}

func boltContext(rootContextID, contextID uint32) proxy.ProtocolContext {
	return &boltProtocolContext{
		bolt:      boltProtocol,
		contextID: contextID,
	}
}

var boltProtocol = bolt.NewBoltProtocol()

type boltProtocolContext struct {
	proxy.DefaultRootContext // notify on plugin start.
	proxy.DefaultProtocolContext
	bolt      proxy.Protocol
	contextID uint32
}

// protocol feature

func (proto *boltProtocolContext) Name() string {
	return proto.bolt.Name()
}

func (proto *boltProtocolContext) Codec() proxy.Codec {
	return proto.bolt.Codec()
}

func (proto *boltProtocolContext) KeepAlive() proxy.KeepAlive {
	return proto.bolt
}

func (proto *boltProtocolContext) Hijacker() proxy.Hijacker {
	return proto.bolt
}

// vm and plugin lifecycle

func (proto *boltProtocolContext) OnVMStart(conf proxy.ConfigMap) bool {

	proxy.Log.Infof("proxy_on_vm_start from Go!, config %v", conf)

	return true
}

func (proto *boltProtocolContext) OnPluginStart(conf proxy.ConfigMap) bool {

	proxy.Log.Infof("proxy_on_plugin_start from Go!")

	return true
}
