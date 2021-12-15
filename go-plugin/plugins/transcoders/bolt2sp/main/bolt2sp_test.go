package main

import (
	"encoding/json"
	"testing"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestBolt2spHttpPath(t *testing.T) {
	type fields struct {
		cfg string
	}
	type args struct {
		headers *bolt.Request
	}

	RequestHeader := bolt.RequestHeader{}
	RequestHeader.Header.Set(ServiceName, "a")
	RequestHeader.Header.Set(BoltMethodName, "b")

	fh := fasthttp.RequestHeader{}
	fh.Set(ServiceName, "a")
	fh.Set(BoltMethodName, "b")
	fh.Set(MosnPath, "/hello")
	fh.Set(ServiceName, "dubbo2http")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   fasthttp.Request
	}{
		{
			name: "TestBolt2spHttpPath",
			fields: fields{
				cfg: `{"service":"dubbo2http","a":{"b":{"x-mosn-path":"/hello"},"c":{"x-mosn-path":"test"}}}`,
			},
			args: args{
				headers: &bolt.Request{
					RequestHeader: RequestHeader,
				},
			},
			want: fasthttp.Request{
				Header: fh,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			cfg := make(map[string]interface{})
			json.Unmarshal([]byte(tt.fields.cfg), &cfg)
			bolt := &bolt2sp{
				cfg: cfg,
			}
			got := bolt.httpReq2BoltReq(tt.args.headers)
			assert.Equal(t, tt.want, got)
		})
	}
}
