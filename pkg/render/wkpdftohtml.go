package render

import (
	"encoding/json"
	"log"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

// GenPDF receives the final path and URL
// URL is passed as is to the wkhtmltopdf
// and outputFilePath is provided as the destination to write file to as is
// meanwhile, object and converter options from mod file are overridden
func (r *Render) GenPDF() (err error) {
	htmlObject, err := pdf.NewObject(r.URL)
	if err != nil {
		return
	}

	// updating html related properties
	htmlObject.UseLocalLinks = true

	// properties from mod file
	// if and only if there is something to override, otherwise json unmarshal gives error
	// e: unexpected end of JSON input
	if len(*r.ObjectOptions) > 0 {
		err = json.Unmarshal(*r.ObjectOptions, &htmlObject)
	}
	if err != nil {
		return
	}

	// wkhtmltopdf converter contains pdf specific configuration
	converter, err := pdf.NewConverter()
	if err != nil {
		return
	}
	defer converter.Destroy()

	// updating pdf related properties
	converter.Title = "Start the fire"
	converter.PaperSize = pdf.A7
	converter.Orientation = pdf.Landscape
	converter.MarginTop = "0mm"
	converter.MarginBottom = "0mm"
	converter.MarginLeft = "0mm"
	converter.MarginRight = "0mm"
	converter.Colorspace = pdf.Grayscale

	// properties from mod file
	// if and only if there is something to override, otherwise json unmarshal gives error
	// e: unexpected end of JSON input
	if len(*r.ConverterOptions) > 0 {
		err = json.Unmarshal(*r.ConverterOptions, &converter)
	}
	if err != nil {
		return
	}

	outFile, err := os.Create(r.OutputFilePath)
	if err != nil {
		return
	}
	defer func() {
		_ = outFile.Close()
		if err != nil {
			err = os.Remove(r.OutputFilePath)
			if err != nil {
				log.Println("Failed to delete corrupt file", err)
			}
		}
	}()

	// render pdf for the above created object
	converter.Add(htmlObject)
	if err = converter.Run(outFile); err != nil {
		return
	}
	return
}
