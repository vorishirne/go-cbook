package render

import (
	"encoding/json"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/velcrine/goreader/pkg/pdfrender"
	"github.com/watergist/file-engine/reader"
	"github.com/watergist/file-engine/webclient"
	"os"
)

// GoReaderRenderTempDir is the tmp dir for goreader
const GoReaderRenderTempDir = "/tmp/goreader/html"

// TempHTMLFilePath if the file type is a .md,  convert to html at this location
const TempHTMLFilePath = GoReaderRenderTempDir + "/current.html"

// TempCSSFilePath all the css from different sources like default and
// overrides will be written here first everytime and then this location
// is set into wkhtmltopdf
const TempCSSFilePath = GoReaderRenderTempDir + "/current.css"

// DefaultCssFile the css that is always there for every host
const DefaultCssFile = "css/default.css"

type Render struct {
	// SiteObjectOptions objectOptions coming from options set in webpages-properties
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	SiteObjectOptions *json.RawMessage
	// SiteConverterOptions objectOptions coming from options set in webpages-properties
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	SiteConverterOptions *json.RawMessage
	// ObjectOptions coming from the mod file
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	ObjectOptions *json.RawMessage
	// ConverterOptions coming from the mod file
	// these options are passed to the wkhtmltopdf converter
	// could be used for overriding default ones in this repository
	ConverterOptions *json.RawMessage
	// URL to webPage or path to html page in local
	URL string
	// OutputFilePath the indexed path to write pdf to
	OutputFilePath string
}

func Do(m *pdfrender.Mod, URL, outputFilePath string) (err error) {
	// set as per in type doc
	render := Render{
		ObjectOptions:        m.ObjectOptions,
		ConverterOptions:     m.ConverterOptions,
		SiteObjectOptions:    m.State.CurrentURLProperties.ObjectOptions,
		SiteConverterOptions: m.State.CurrentURLProperties.ConverterOptions,
		URL:                  URL,
		OutputFilePath:       outputFilePath,
	}
	// ensure temp directory is created
	err = os.MkdirAll(GoReaderRenderTempDir, 0766)
	if err != nil {
		return
	}
	if m.State.IsMD {
		err = render.CompileMDtoHTML()
		if err != nil {
			return
		}
	}
	// all the css files getting piled up to one css file
	err = reader.CopyContent(TempCSSFilePath, DefaultCssFile,
		m.State.CurrentURLProperties.CssOverrideFile, m.CssOverrideFile)
	if err != nil {
		return
	}

	// finally, render the pdf
	err = render.GenPDF()
	return
}

func (r *Render) CompileMDtoHTML() (err error) {
	// get the bytes from http from this url
	md, err := webclient.HTTPResponse(r.URL)
	// package syntax
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(*md, mdParser, nil)

	// write the markdown to a temp file and set that file as current html
	err = os.WriteFile(r.URL, html, 0666)
	r.URL = TempHTMLFilePath
	return
}
