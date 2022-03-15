package encryption

import (
	"reflect"
	"testing"
)

func TestXorEncoder(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		msg  []byte
		want []byte
	}{
		{
			name: "XorEncoder",
			key:  []byte{7, 2, 34, 9, 10, 6, 78},
			msg:  []byte{1, 5, 8, 3, 5, 24, 89, 0, 35, 10, 39, 127, 46, 97, 11, 16, 18, 49, 29, 47, 14, 29},
			want: []byte{6, 7, 42, 10, 15, 30, 23, 7, 33, 40, 46, 117, 40, 47, 12, 18, 48, 56, 23, 41, 64, 26},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := XorEncrypt(tt.msg, tt.key)
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("XorEncoder() got = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestXorDecoder(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		msg  []byte
		want []byte
	}{
		{
			name: "XorDecoder",
			key:  []byte{7, 2, 34, 9, 10, 6, 78},
			msg:  []byte{6, 7, 42, 10, 15, 30, 23, 7, 33, 40, 46, 117, 40, 47, 12, 18, 48, 56, 23, 41, 64, 26},
			want: []byte{1, 5, 8, 3, 5, 24, 89, 0, 35, 10, 39, 127, 46, 97, 11, 16, 18, 49, 29, 47, 14, 29},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := XorDecrypt(tt.msg, tt.key)
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("XorDecoder() got = %v, want %v", res, tt.want)
			}
		})
	}
}
