package generate

import (
	"encoding/json"
	"os"
	"path"
)

type DirVisited struct {
	IndexedDirPath string
	ItemCount      int8
	DirIndex       string
}
type Mod struct {
	BaseUrl          string
	BaseDir          string
	HistPointer      string
	ObjectOptions    json.RawMessage
	ConverterOptions json.RawMessage
	dirVisited       map[string]*DirVisited
	baseDirLastIndex int8
}

func AllPagesMod(modFilePath string) (m *Mod, err error) {
	modFile, err := os.ReadFile(modFilePath)
	if err != nil {
		return
	}
	m = &Mod{}
	m.dirVisited = map[string]*DirVisited{}
	err = json.Unmarshal(modFile, m)
	return
}

func (m *Mod) GetFilePath(webUrl string) (filePath string, err error) {
	rawDirPath, fileName, err := m.GetRawFilePath(webUrl)
	if err != nil {
		return
	}
	parentDir, itemIndex, err := m.GetIndexedDir(rawDirPath)
	if err != nil {
		return
	}
	err = os.MkdirAll(parentDir, 0766)
	if err != nil {
		return
	}
	filePath = path.Join(parentDir, itemIndex+fileName+".pdf")
	return
}
