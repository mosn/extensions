// How to use conditional compilation with the go build tool:
// https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool

// +build proxytest

package proxy

import (
	"reflect"
	"unsafe"
)

// parseString parse byte pointer to string
func parseString(buf *byte, len int) string {
	if len <= 0 || buf == nil {
		return ""
	}

	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buf)),
		Len:  len,
		Cap:  len,
	}))
}

// parseByteSlice parse byte pointer to byte slice
func parseByteSlice(buf *byte, len int) []byte {
	if len <= 0 || buf == nil {
		return []byte{}
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buf)),
		Len:  len,
		Cap:  len,
	}))
}
