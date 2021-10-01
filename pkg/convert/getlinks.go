package convert

import (
	"encoding/json"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"log"
	"os"
)

type DirVisited struct {
	DirPath           string
	CountReached      int8
	LastCounterString string
}
type Mod struct {
	BaseUrl             string
	BaseDir             string
	HistPointer         string
	ObjectOptions       json.RawMessage
	ConverterOptions    json.RawMessage
	dirVisited          map[string]*DirVisited
	baseDirCountReached int8
}

func AllPagesMod(modFilePath string) (m *Mod, err error) {
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	modFile, err := os.ReadFile(modFilePath)
	if err != nil {
		return
	}
	m = &Mod{}
	m.dirVisited = map[string]*DirVisited{}
	err = json.Unmarshal(modFile, m)
	return
}
