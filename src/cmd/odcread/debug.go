//go:build debug
// +build debug

package main

import (
	"fmt"
	"os"
)

func debugLog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
}
