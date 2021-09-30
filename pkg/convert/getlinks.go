package convert

import (
	"encoding/json"
	"fmt"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

type DirVisited struct {
	DirPath           string
	CountReached      int8
	LastCounterString string
}
type Mod struct {
	CssFile             string
	BaseUrl             string
	BaseDir             string
	HistPointer         uint8
	ObjectOptions       json.RawMessage
	ConverterOptions    json.RawMessage
	Padding             map[string]int8
	dirVisited          map[string]*DirVisited
	BaseDirCountReached int8
}

func AllPages(urlsFilePath, modFilePath string) {
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	//writeConverted("https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request",
	//	"../kindle-web/envoy/css.html", "lifeofreq", false)
	//writeConverted("https://blog.envoyproxy.io/how-to-write-envoy-filters-like-a-ninja-part-1-d166e5abec09",
	//	"css/medium.css", "filter2", true)
	//writeConverted("https://blog.envoyproxy.io/taming-a-network-filter-44adcf91517",
	//	"css/medium.css", "filter3", true)
	writeConverted("https://medium.com/@abhbose6/bazel-101-2b0272b15da8",
		"css/medium.css", "bazel2", true)
	writeConverted("https://medium.com/@alishananda/implementing-filters-in-envoy-dcd8fc2d8fda",
		"css/medium.css", "filter4", true)
}

func (m *Mod) GetFilePath(webUrl string) (dirPath, fileName string, err error) {

	webUrlObj, err := url.Parse(webUrl)
	if err != nil {
		return "", "", err
	}

	baseUrlObj, err := url.Parse(m.BaseUrl)
	if err != nil {
		return "", "", err
	}
	webUrlPath := strings.Trim(webUrlObj.Path, "/")
	baseUrlPath := strings.Trim(baseUrlObj.Path, "/")
	if !strings.HasPrefix(webUrlPath, baseUrlPath) {
		return "", "", fmt.Errorf("mod base url %v not valid for %v", webUrlPath, baseUrlPath)
	}

	dirPath = strings.Trim(strings.TrimPrefix(webUrlPath, baseUrlPath), "/")
	if dirPath == "" {
		return "", "", fmt.Errorf("filepath is empty")
	}
	lastSlash := strings.LastIndexAny(dirPath, "/")

	fileName = dirPath[lastSlash+1:]
	if lastSlash == -1 {
		lastSlash++
	}
	dirPath = dirPath[:lastSlash]
	if dirPath == "" {
		return "", "", fmt.Errorf("no directory to write file to")
	}
	return
}

func (m *Mod) EnsureNumberedDir(dirPath string) (parentDir, fileCounter string, err error) {
	if info, ok := m.dirVisited[dirPath]; ok {
		info.CountReached++
		return info.DirPath, info.LastCounterString + "-" + fmt.Sprintf("%03d", info.CountReached), nil
	}
	lastSlash := strings.LastIndexAny(dirPath, "/")
	if lastSlash == -1 {
		m.BaseDirCountReached++
		newFileCounter := fmt.Sprintf("%03d", m.BaseDirCountReached)
		newDirName := path.Join(m.BaseDir, newFileCounter+"-"+dirPath)

		//common code
		err = os.MkdirAll(newDirName, 0766)
		if err != nil {
			return
		}
		m.dirVisited[dirPath] = &DirVisited{
			newDirName, 1, newFileCounter}

		return newDirName, newFileCounter + "-" + fmt.Sprintf("%03d", 1), nil
	}
	suffixDirPath := dirPath[lastSlash+1:]
	prefixDirPath := dirPath[:lastSlash]
	newDirName, newFileCounter, err := m.EnsureNumberedDir(prefixDirPath)
	if err != nil {
		return
	}
	parentDir = path.Join(newDirName, newFileCounter+"-"+suffixDirPath)
	fileCounter = newFileCounter + "-" + fmt.Sprintf("%03d", 1)

	// common code
	err = os.Mkdir(parentDir, 0766)
	if err != nil {
		return
	}
	m.dirVisited[dirPath] = &DirVisited{
		parentDir, 1, newFileCounter}

	return
}

func createAndUpdateDirInfo(dirName, fileCounter string) {

}
