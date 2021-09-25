package convert

import (
	"log"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

func writeConverted(URL, cssFile, outputFile string, disableJS bool) {

	object2, err := pdf.NewObject(URL)
	object2.UserStylesheetLocation = cssFile
	object2.UseLocalLinks = true
	object2.EnableJavascript = !disableJS

	if err != nil {
		log.Fatal(err)
	}

	// Create object from reader.
	converter, err := pdf.NewConverter()

	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add created objects to the converter.
	converter.Add(object2)

	// Set converter options.
	converter.Title = "Sample document"
	converter.PaperSize = pdf.B8
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "4mm"
	converter.MarginBottom = "4mm"
	converter.MarginLeft = "2mm"
	converter.MarginRight = "2mm"
	converter.Colorspace = pdf.Grayscale

	// Convert objects and save the output PDF document.
	outFile, err := os.Create(outputFile + ".pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	if err := converter.Run(outFile); err != nil {
		log.Fatal(err)
	}
}
