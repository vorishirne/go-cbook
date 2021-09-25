package convert

import (
	pdf "github.com/adrg/go-wkhtmltopdf"
	"log"
)

type Mod struct {
	CssFile            string
	BasePath           string
	BaseDirReplaceName string
	HistPointer        uint8
	DisableJS          bool
}

func AllPages(urlsFilePath, modFilePath string) {
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	writeConverted("https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request",
		"../kindle-web/envoy/css.html", "lifeofreq", false)
	writeConverted("https://blog.envoyproxy.io/external-c-dependency-management-in-bazel-dd37477422f5",
		"css/medium.css", "filter1", true)
	writeConverted("https://blog.envoyproxy.io/how-to-write-envoy-filters-like-a-ninja-part-1-d166e5abec09",
		"css/medium.css", "filter2", true)
	writeConverted("https://blog.envoyproxy.io/taming-a-network-filter-44adcf91517",
		"css/medium.css", "filter3", true)

}

func ClearUntrackedHist(modFilePath string) {

}
