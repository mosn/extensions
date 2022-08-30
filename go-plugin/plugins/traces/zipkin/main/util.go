package main

import (
	"strings"
)

func getOperationName(uri []byte) string {
	arr := strings.Split(string(uri), "?")
	return arr[0]
}
