package pdfrender

import (
	"encoding/json"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

// GenPDF receives the final path and URL
// URL is passed as is to the wkhtmltopdf
// and outputFilePath is provided as the destination to write file to as is
// meanwhile, object and converter options from mod file are overridden
func (m *Mod) GenPDF(URL, outputFilePath string) (err error) {
	// wkhtmltopdf object contains html specific configuration
	object, err := pdf.NewObject(URL)
	if err != nil {
		return
	}

	// updating html related properties
	object.UseLocalLinks = true
	// properties from mod file
	err = json.Unmarshal(m.ObjectOptions, &object)
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
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1mm"
	converter.MarginBottom = "0mm"
	converter.MarginLeft = "1mm"
	converter.MarginRight = "0mm"
	converter.Colorspace = pdf.Grayscale
	// properties from mod file
	err = json.Unmarshal(m.ConverterOptions, &converter)
	if err != nil {
		return
	}

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return
	}
	defer func() {
		_ = outFile.Close()
	}()

	// render pdf for the above object
	converter.Add(object)
	if err = converter.Run(outFile); err != nil {
		return
	}
	return
}