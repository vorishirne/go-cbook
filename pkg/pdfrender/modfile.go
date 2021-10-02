package pdfrender

import (
	"encoding/json"
	"os"
	"path"
)

// DirVisited allows you to save below information for multiple urls
// not just for caching purpose, but also for to keep track of info for created dirs
type DirVisited struct {
	// the indexed path for this dir
	IndexedDirPath string
	// the no. of items currently in this dir
	ItemCount int8
	// the complex index like i01-03-01 is carried from parent directory,
	// that needs to be available for sending to create the item in this dir.
	DirIndex string
}
type Mod struct {
	// the base url to be taken away as prefix from the urls' paths
	BaseUrl string
	// base parent dir to all the dirs
	BaseDir string
	// a rune, to mark the current running set of files & dirs, that will be
	// differentiated from the past available files &dirs
	HistPointer string
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	ObjectOptions json.RawMessage
	// these options are passed to the wkhtmltopdf converter
	// could be used for overriding default ones in this repository
	ConverterOptions json.RawMessage
	//				 map[rawFilePath]dataForIndex
	dirVisited map[string]*DirVisited
	// no. of items in Mod.BaseDir
	baseDirLastIndex int8
}

// GetMod prepares the Mod struct from mod json file provided
func GetMod(modFilePath string) (m *Mod, err error) {
	modFile, err := os.ReadFile(modFilePath)
	if err != nil {
		return
	}
	m = &Mod{}
	// dirVisited is hidden and not expected to come from json
	m.dirVisited = map[string]*DirVisited{}
	err = json.Unmarshal(modFile, m)
	return
}

// GetFilePath wraps away the logic to provide & create
// indexed dir from the url path or url fragment
func (m *Mod) GetFilePath(webUrl string) (filePath string, err error) {
	// 1. Extracts raw filepath from url.
	rawDirPath, fileName, err := m.GetRawFilePath(webUrl)
	if err != nil {
		return
	}
	// 2. Converts it to indexed one
	parentDir, itemIndex, err := m.GetIndexedDir(rawDirPath)
	if err != nil {
		return
	}
	//  & creates that dir, even if its already there
	err = os.MkdirAll(parentDir, 0766)
	if err != nil {
		return
	}
	// 3. Returns the final filePath to be used to create pdf at.
	filePath = path.Join(parentDir, itemIndex+fileName+".pdf")
	return
}
