package render

import (
	"encoding/json"
	"log"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

// GenPDF receives the final render object
// OutputFilePath is provided as the destination to write file to as is
// Also, object and converter options from mod file and webpage-properties
// files do override defaultly set
func (r *Render) GenPDF() (err error) {
	htmlObject, err := pdf.NewObject(r.URL)
	if err != nil {
		return
	}

	// updating html related properties
	htmlObject.UseLocalLinks = true
	htmlObject.UserStylesheetLocation = TempCSSFilePath
	// if and only if there is something to override, otherwise json unmarshal gives error
	// e: unexpected end of JSON input

	// options from webpage-properties file
	if r.SiteObjectOptions != nil && len(*r.SiteObjectOptions) > 0 {
		err = json.Unmarshal(*r.SiteObjectOptions, &htmlObject)
	}
	if err != nil {
		return
	}
	// options from mod file
	if r.ObjectOptions != nil && len(*r.ObjectOptions) > 0 {
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
	// https://doc.qt.io/archives/qt-4.8/qprinter.html#PaperSize-enum
	converter.Height = "76"
	converter.Width = "100"

	converter.Orientation = pdf.Portrait
	converter.MarginTop = "0mm"
	converter.MarginBottom = "0mm"
	converter.MarginLeft = "0mm"
	converter.MarginRight = "0mm"
	converter.Colorspace = pdf.Grayscale

	// if and only if there is something to override, otherwise json unmarshal gives error
	// e: unexpected end of JSON input

	// options from webpage-properties file
	if r.SiteConverterOptions != nil && len(*r.SiteConverterOptions) > 0 {
		err = json.Unmarshal(*r.SiteConverterOptions, &converter)
	}
	if err != nil {
		return
	}
	// options from mod file
	if r.ConverterOptions != nil && len(*r.ConverterOptions) > 0 {
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
		// if the response is an error, delete the created empty file
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
