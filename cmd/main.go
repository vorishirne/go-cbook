package main

import (
	"fmt"
	"github.com/velcrine/eink-pages/pkg/convert"
	"log"
	"runtime"
)

func init() {
	// Set main function to run on the main thread.
	runtime.LockOSThread()
}

func main() {
	target := "istio-docs"
	err := convert.ReadUrlFile(fmt.Sprintf("urls/%v.url", target),
		fmt.Sprintf("mods/%v.json", target))
	if err != nil {
		log.Fatal(err)
	}
}
