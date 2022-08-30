package main

import "log"

// NewDefaultLogger Creates a new Log
func NewDefaultLogger() *log.Logger {
	return log.Default()
}
