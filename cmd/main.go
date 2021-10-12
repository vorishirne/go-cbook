package main

import (
	"fmt"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/velcrine/eink-pages/pkg/config"
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

	// flags, glob, global css, readme, ensureDir fn
	// repeat for file-engine and
	// committer : see difference
	// release both
	// git config global in all
	// reach upto mohit in galgodas, remv mohit professor name

	// make galgodas ready: as to remove main package to respective, and remove indexing
	// import flutter repo, as it has nothing as part of improved commit. : change commit time, to september, so as to make the commits count for long

	// one more constraint from qt,
	// init must be called from main function

	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

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
}
