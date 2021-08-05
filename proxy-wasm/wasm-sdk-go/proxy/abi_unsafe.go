package proxy

import "unsafe"

// parseBytes parse string to byte pointer
func stringBytePtr(message string) *byte {
	if len(message) == 0 {
		return nil
	}

	buffer := *(*[]byte)(unsafe.Pointer(&message))
	return &buffer[0]
}

func parseSliceString(buf []byte) string {
	if len(buf) <= 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&buf))
}
