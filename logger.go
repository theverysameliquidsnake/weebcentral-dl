package main

import (
	"time"
	"fmt"
	"log"
)

type customWriter struct{}

func (cw *customWriter) Write(bs []byte) (int, error) {
	return fmt.Print("[", time.Now().Format("15:04:05"), "] ", string(bs))
}

func debugOutput(msg string) {
	if isDebugOutputEnabled {
		log.Println("[DEBUG]", msg)
	}
}

func concatErrorString(prefix string, err error) string {
	return fmt.Sprintf(prefix, err)
}
