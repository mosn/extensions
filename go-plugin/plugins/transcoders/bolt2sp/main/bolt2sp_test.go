package main

import (
	"encoding/json"
	"testing"

	"github.com/mosn/extensions/go-plugin/pkg/protocol/bolt"
	"github.com/stretchr/testify/assert"
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
	RequestHeader.Header.Set(MethodName, "b")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "TestBolt2spHttpPath",
			fields: fields{
				cfg: `{"a":{"b":"/hello","c":"test"}}`,
			},
			args: args{
				headers: &bolt.Request{
					RequestHeader: RequestHeader,
				},
			},
			want: "/hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			cfg := make(map[string]interface{})
			json.Unmarshal([]byte(tt.fields.cfg), &cfg)
			bolt := &bolt2sp{
				cfg: cfg,
			}
			got := bolt.httpPath(tt.args.headers)
			assert.Equal(t, tt.want, got)
		})
	}
}
