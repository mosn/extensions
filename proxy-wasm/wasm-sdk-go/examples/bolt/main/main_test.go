package main

import (
	"context"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/examples/bolt"
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy"
	"reflect"
	"testing"
)

func TestBolt(t *testing.T) {

	vmConfig := proxy.NewConfigMap()
	vmConfig.Set("engine", "wasm")

	opt := proxy.NewEmulatorOption().
		WithNewProtocolContext(boltContext).
		WithNewRootContext(rootContext).
		WithVMConfiguration(vmConfig)

	host := proxy.NewHostEmulator(opt)
	// release lock and reset emulator state
	defer host.Done()
	// invoke host start vm
	host.StartVM()
	// invoke host plugin
	host.StartPlugin()

	// 1. invoke downstream decode
	ctxId := host.NewProtocolContext()
	// bolt plugin decode will be invoked
	cmd, err := host.Decode(ctxId, proxy.WrapBuffer(decodedRequestBytes(uint32(host.CurrentStreamId()))))
	if err != nil {
		t.Fatalf("failed to invoke host decode request buffer, err: %v", err)
	}

	if _, ok := cmd.(*bolt.Request); !ok {
		t.Fatalf("decode request type error, expect *bolt.Request, actual %v", reflect.TypeOf(cmd))
	}

	// 2. invoke upstream encode
	upstreamBuf, err := host.Encode(ctxId, cmd)
	if err != nil {
		t.Fatalf("failed to invoke host encode request buffer, err: %v", err)
	}

	// check upstream content with downstream request
	if !reflect.DeepEqual(decodedRequestBytes(uint32(host.CurrentStreamId())), upstreamBuf.Bytes()) {
		t.Fatalf("failed to invoke host encode request buffer, err: %v", err)
	}

	// complete protocol pipeline
	host.CompleteProtocolContext(ctxId)
}

func TestBoltHijack(t *testing.T) {

	vmConfig := proxy.NewConfigMap()
	vmConfig.Set("engine", "wasm")

	opt := proxy.NewEmulatorOption().
		WithNewProtocolContext(boltContext).
		WithNewRootContext(rootContext).
		WithVMConfiguration(vmConfig)

	host := proxy.NewHostEmulator(opt)
	// release lock and reset emulator state
	defer host.Done()
	// invoke host start vm
	host.StartVM()
	// invoke host plugin
	host.StartPlugin()

	// 1. invoke downstream decode
	ctxId := host.NewProtocolContext()
	// bolt plugin decode will be invoked
	cmd, err := host.Decode(ctxId, proxy.WrapBuffer(decodedRequestBytes(uint32(host.CurrentStreamId()))))
	if err != nil {
		t.Fatalf("failed to invoke host decode request buffer, err: %v", err)
	}

	if _, ok := cmd.(*bolt.Request); !ok {
		t.Fatalf("decode request type error, expect *bolt.Request, actual %v", reflect.TypeOf(cmd))
	}

	// 2. invoke upstream encode
	upstreamBuf, err := host.Encode(ctxId, cmd)
	if err != nil {
		t.Fatalf("failed to invoke host encode request buffer, err: %v", err)
	}

	// check upstream content with downstream request
	if !reflect.DeepEqual(decodedRequestBytes(uint32(host.CurrentStreamId())), upstreamBuf.Bytes()) {
		t.Fatalf("failed to invoke host encode request buffer, err: %v", err)
	}

	// complete protocol pipeline
	host.CompleteProtocolContext(ctxId)

	// 3. mock request failed, hijack triggered
	hijackId := host.NewProtocolContext()
	resp := host.Hijack(hijackId, cmd.(*bolt.Request), 504)
	_ = resp
	host.CompleteProtocolContext(hijackId)
}

func decodedRequestBytes(id uint32) []byte {
	rpcHeader := proxy.NewHeader()
	rpcHeader.Set("service", "com.alipay.demo.HelloService")
	request := bolt.NewRpcRequest(id, rpcHeader, proxy.WrapBuffer([]byte("bolt body")))
	buf, err := bolt.NewBoltProtocol().Codec().Encode(context.TODO(), request)
	if err != nil {
		panic("failed to encode bolt request, err: " + err.Error())
	}
	return buf.Bytes()
}
