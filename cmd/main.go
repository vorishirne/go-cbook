package main

import (
	"fmt"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/velcrine/goreader/pkg/config"
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
	GeneratePDFs()
}

func GeneratePDFs() {

	// this is the name of mod file & url file
	target := "istio-docs"

	// RenderFromUrlFile is the root function to iterate
	// over every url present in the provided url file
	// and calls pdf-render process
	err := config.RenderFromUrlFile(fmt.Sprintf("urls/%v.url", target),
		fmt.Sprintf("mods/%v.json", target))
	if err != nil {
		log.Fatal(err)
	}

	target = "ssh-keys"

	// RenderFromUrlFile is the root function to iterate
	// over every url present in the provided url file
	// and calls pdf-render process
	err = config.RenderFromUrlFile(fmt.Sprintf("urls/%v.url", target),
		fmt.Sprintf("mods/%v.json", target))
	if err != nil {
		log.Fatal(err)
	}

	target = "go-blogs"

	// RenderFromUrlFile is the root function to iterate
	// over every url present in the provided url file
	// and calls pdf-render process
	err = config.RenderFromUrlFile(fmt.Sprintf("urls/%v.url", target),
		fmt.Sprintf("mods/%v.json", target))
	if err != nil {
		log.Fatal(err)
	}
}
