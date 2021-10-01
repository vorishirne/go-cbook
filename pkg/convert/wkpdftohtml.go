package convert

import (
	"encoding/json"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

func (m *Mod) GenPDF(URL, outputFilePath string) (err error) {

	object, err := pdf.NewObject(URL)
	if err != nil {
		return
	}

	object.UseLocalLinks = true
	err = json.Unmarshal(m.ObjectOptions, &object)
	if err != nil {
		return
	}

	converter, err := pdf.NewConverter()
	if err != nil {
		return
	}
	defer converter.Destroy()

	converter.Title = "Start the fire"
	converter.PaperSize = pdf.A7
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1mm"
	converter.MarginBottom = "0mm"
	converter.MarginLeft = "1mm"
	converter.MarginRight = "0mm"
	converter.Colorspace = pdf.Grayscale
	err = json.Unmarshal(m.ConverterOptions, &converter)

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return
	}
	defer func() {
		_ = outFile.Close()
	}()

	converter.Add(object)
	if err = converter.Run(outFile); err != nil {
		return
	}
	return
}
