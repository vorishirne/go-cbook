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

const GoReaderRenderTempDir = "/tmp/goreader/html"
const TempHTMLFilePath = GoReaderRenderTempDir + "/current.html"
const TempCSSFilePath = GoReaderRenderTempDir + "/current.css"
const DefaultCssFile = "css/default.css"

type Render struct {
	SiteObjectOptions    *json.RawMessage
	SiteConverterOptions *json.RawMessage
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	ObjectOptions *json.RawMessage
	// these options are passed to the wkhtmltopdf converter
	// could be used for overriding default ones in this repository
	ConverterOptions *json.RawMessage
	// the url to webPage or path to html page in local
	URL string

	OutputFilePath string

	CSSOverridePath string
}

func Do(m *pdfrender.Mod, URL, outputFilePath string) (err error) {
	render := Render{
		ObjectOptions:        m.ObjectOptions,
		ConverterOptions:     m.ConverterOptions,
		SiteObjectOptions:    m.State.CurrentURLProperties.ObjectOptions,
		SiteConverterOptions: m.State.CurrentURLProperties.ConverterOptions,
		URL:                  URL,
		OutputFilePath:       outputFilePath,
		CSSOverridePath:      TempCSSFilePath,
	}
	if m.IsMD {
		err = render.CompileMDtoHTML()
		if err != nil {
			return
		}
	}

	err = reader.CopyContent(TempCSSFilePath, DefaultCssFile,
		m.State.CurrentURLProperties.CssOverrideFile, m.CssOverrideFile)
	if err != nil {
		return
	}

	err = render.GenPDF()
	return
}

func (r *Render) CompileMDtoHTML() (err error) {
	md, err := webclient.HTTPResponse(r.URL)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)

	html := markdown.ToHTML(*md, mdParser, nil)

	err = os.MkdirAll(GoReaderRenderTempDir, 0766)
	if err != nil {
		return
	}
	err = os.WriteFile(r.URL, html, 0666)
	r.URL = TempHTMLFilePath
	return
}
