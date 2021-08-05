package proxy

import (
	"fmt"
	"testing"
)

func TestCopyString(t *testing.T) {

	helloWorld := "hello world!"
	buf := make([]byte, len(helloWorld))
	copy(buf[:], helloWorld)
	fmt.Printf("buf => %s\n", string(buf))

}

func TestDecodeMap(t *testing.T) {
	maps := make(map[string]string, 2)
	maps["key1"] = "value1"
	maps["hello"] = "world"

	bytes := EncodeMap(maps)
	decoded := DecodeMap(bytes)

	if val, ok := decoded["key1"]; !ok || val != "value1" {
		t.Errorf("expect value 'value1' for key 'key1', actual '%s'", val)
	}

	if val, ok := decoded["hello"]; !ok || val != "world" {
		t.Errorf("expect value 'world' for key 'hello', actual '%s'", val)
	}
}
