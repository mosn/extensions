package main

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
	"mosn.io/api"
)

func TestBums2CdIterm_GetXml(t *testing.T) {
	type fields struct {
		Key    string
		Type   string
		Length string
		Reader BumsReader
		Scale  string
	}
	type args struct {
		header   api.HeaderMap
		headBody *fastjson.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *etree.Element
		want1  bool
	}{
		{
			fields: fields{
				Key:    "Key",
				Type:   "string",
				Length: "8",
				Reader: BumsReader{
					Origin:  "key",
					Key:     "key",
					Default: "default",
				},
				Scale: "0",
			},
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iterm := Bums2CdIterm{
				Key:    tt.fields.Key,
				Type:   tt.fields.Type,
				Length: tt.fields.Length,
				Reader: tt.fields.Reader,
				Scale:  tt.fields.Scale,
			}
			got, got1 := iterm.GetXml(tt.args.header, tt.args.headBody)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, got1, tt.want1)
		})
	}
}
