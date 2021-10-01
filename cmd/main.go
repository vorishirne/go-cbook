package main

import (
	"fmt"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/velcrine/eink-pages/pkg/convert"
	"log"
	"runtime"
)

func init() {
	// Set main function to run on the main thread.
	runtime.LockOSThread()
}

func main() {
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	target := "istio-docs"
	err := convert.ReadUrlFile(fmt.Sprintf("urls/%v.url", target),
		fmt.Sprintf("mods/%v.json", target))
	if err != nil {
		log.Fatal(err)
	}
}
