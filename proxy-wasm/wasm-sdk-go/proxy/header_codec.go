package proxy

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

var (
	errInvalidLength = errors.New("invalid length -1, ignore current key value pair")
)

func GetEncodeHeaderLength(h Header) int {
	n := 0
	h.Range(func(key, value string) bool {
		n += 8 + len(key) + len(value)
		return true
	})
	return n
}

func EncodeHeader(buf Buffer, h Header) {
	h.Range(func(key, value string) bool {
		encodeString(buf, key)
		encodeString(buf, value)
		return true
	})
}

func DecodeHeader(bytes []byte, h Header) (err error) {
	totalLen := len(bytes)
	index := 0
	var key, value []byte

	for index < totalLen {
		// 1. read key
		key, index, err = decodeString(bytes, totalLen, index)
		if err != nil {
			if err == errInvalidLength {
				continue
			}
			return
		}

		// 2. read value
		value, index, err = decodeString(bytes, totalLen, index)
		if err != nil {
			if err == errInvalidLength {
				continue
			}
			return
		}

		// 3. kv append
		h.Set(string(key), string(value))
	}
	return nil
}

func encodeString(buf Buffer, str string) {
	length := len(str)
	// 1. encode str length
	buf.WriteUint32(uint32(length))
	// 2. encode str value
	buf.Write([]byte(str))
}

func decodeString(bytes []byte, totalLen, index int) (str []byte, newIndex int, err error) {
	// 1. read str length
	length := binary.BigEndian.Uint32(bytes[index:])

	// avoid length = -1
	if length == math.MaxUint32 {
		return nil, index + 4, errInvalidLength
	}

	end := index + 4 + int(length)
	if end > totalLen {
		return nil, end, fmt.Errorf("decode header failed, index %d, length %d, totalLen %d, bytes %v\n", index, length, totalLen, bytes)
	}

	return bytes[index+4 : end], end, nil
}
