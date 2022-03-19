package encryption

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
)

func TestParseSecret(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  *SecretConfig
	}{
		{
			name:  "ParseSecret",
			value: "{\"enable\":true, \"type\":\"xor\", \"secrets\":{\"ESB002\":\"12345678\",\"ESB001\":\"13213211\"}}",
			want: &SecretConfig{
				Enable: true,
				Type:   "xor",
				Secret: map[string]string{
					"ESB002": "12345678",
					"ESB001": "13213211",
				},
			},
		},
		{
			name:  "ParseSecret",
			value: "{\"enable\":false, \"type\":\"xor\", \"secrets\":{\"ESB002\":\"12345678\",\"ESB001\":\"13213211\"}}",
			want:  nil,
		},
		{
			name:  "ParseSecret",
			value: "",
			want:  nil,
		},
		{
			name:  "ParseSecret",
			value: "{\"type\":\"xor\", \"secrets\":{\"ESB002\":\"12345678\",\"ESB001\":\"13213211\"}}",
			want:  nil,
		},
		{
			name:  "ParseSecret",
			value: "{cdc}",
			want:  nil,
		},
	}

	for _, tt := range tests {
		atoValue := &atomic.Value{}
		atoValue.Store(tt.value)
		ctx := context.WithValue(context.Background(), "code_config", atoValue)
		got, _ := ParseSecret(ctx)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseSecret() got = %v, want %v", got, tt.want)
		}
	}
}
