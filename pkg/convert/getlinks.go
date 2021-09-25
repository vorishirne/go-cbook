package convert

type Mod struct {
	CssFile string
	BasePath string
	BaseDirReplaceName string
	HistPointer uint8
}
func AllPages(urlsFilePath,modFilePath string)  {
	writeConverted("https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request",
		"css/envoy.css","lifeofreq")
	writeConverted("https://blog.envoyproxy.io/external-c-dependency-management-in-bazel-dd37477422f5",
		"css/medium.css","filter1")
	writeConverted("https://blog.envoyproxy.io/taming-a-network-filter-44adcf91517",
		"css/medium.css","filter2")

}

func ClearUntrackedHist(modFilePath string){

}