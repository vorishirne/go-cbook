package main

import (
	"fmt"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/velcrine/eink-pages/pkg/pdfrender"
	"log"
	"runtime"
)

func init() {
	// Set main function to run on the main thread.
	// this is required by wkhtmltopdf, as it depends on QT
	// which creates this constraint when called from c api
	runtime.LockOSThread()
}

func main() {
	// one more constraint from qt,
	// init must be called from main function

	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	// this is the name of mod file & url file
	target := "go-blogs"

	// RenderFromUrlFile is the root function to iterate
	// over every url present in the provided url file
	// and calls pdfrender process
	err := pdfrender.RenderFromUrlFile(fmt.Sprintf("urls/%v.url", target),
		fmt.Sprintf("mods/%v.json", target))
	if err != nil {
		log.Fatal(err)
	}
}
