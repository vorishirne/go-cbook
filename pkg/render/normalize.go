package render

import (
	"encoding/json"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/velcrine/eink-pages/pkg/pdfrender"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

const GoReaderRenderTempDir = "/tmp/goreader/html"
const TempHTMLFilePath = GoReaderRenderTempDir + "/current.html"

type Render struct {
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	ObjectOptions json.RawMessage
	// these options are passed to the wkhtmltopdf converter
	// could be used for overriding default ones in this repository
	ConverterOptions json.RawMessage
	// the url to webPage or path to html page in local
	URL string

	OutputFilePath string
}

func Do(m *pdfrender.Mod, URL, outputFilePath string) (err error) {
	render := Render{
		ObjectOptions:    m.ObjectOptions,
		ConverterOptions: m.ConverterOptions,
		URL:              URL,
		OutputFilePath:   outputFilePath,
	}
	if m.IsMD {
		err = render.CompileMDtoHTML()
	}
	if err != nil {
		return
	}
	//render.GenPDF
	return
}

func (r *Render) CompileMDtoHTML() (err error) {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	response, _ := netClient.Get(r.URL)
	md, err := ioutil.ReadAll(response.Body)

	// continue with markdown processing
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)

	html := markdown.ToHTML(md, mdParser, nil)

	err = os.MkdirAll(GoReaderRenderTempDir, 0766)
	if err != nil {
		return
	}
	r.URL = TempHTMLFilePath
	err = os.WriteFile(r.URL, html, 0666)
	return
}
