package main

import (
	"github.com/velcrine/eink-pages/pkg/convert"
	"runtime"
)

func init() {
	// Set main function to run on the main thread.
	runtime.LockOSThread()
}

func main() {
	convert.AllPages(",", "")
}
