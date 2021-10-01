package convert

import (
	"fmt"
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
	fmt.Print("converter")
	converter, err := pdf.NewConverter()

	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add created objects to the converter.
	converter.Add(object2)

	// Set converter options.
	converter.Title = "Start the fire"
	converter.PaperSize = pdf.A7
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1mm"
	converter.MarginBottom = "0mm"
	converter.MarginLeft = "1mm"
	converter.MarginRight = "0mm"
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
