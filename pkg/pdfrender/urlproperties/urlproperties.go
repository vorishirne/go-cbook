package urlproperties

import (
	"encoding/json"
	"fmt"
	"github.com/velcrine/goreader/pkg/pdfrender"
	"github.com/watergist/file-engine/reader/structure"
	"math"
	"strings"
)

const DefaultPropertiesFile = "mods/webpages-properties.json"

type URLProperties struct {
	Path string
	// this file is merged with global css rules and passed to wkhtmltopdf
	CssOverrideFile string
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	ObjectOptions *json.RawMessage
	// these options are passed to the wkhtmltopdf converter
	// could be used for overriding default ones in this repository
	ConverterOptions *json.RawMessage
}

var allPropertiesFile map[string]*URLProperties

func Init(propertiesFiles []string) (err error) {
	allPropertiesFile = make(map[string]*URLProperties)
	err, _ = structure.LoadJsonFile(DefaultPropertiesFile,
		&allPropertiesFile)
	if err != nil {
		return
	}
	for _, v := range propertiesFiles {
		err, _ = structure.LoadJsonFile(v,
			&allPropertiesFile)
		if err != nil {
			return
		}
	}
	return
}

func SetURLProperties(m *pdfrender.Mod, url string) (properties *URLProperties, err error) {
	err = Init(m.PropertiesFiles)
	if err != nil {
		return
	}
	lowestIndex := math.MaxInt64
	key := ""
	for k := range allPropertiesFile {
		if index := strings.Index(url, k); index > -1 {
			if index < lowestIndex {
				key = k
				lowestIndex = index
			}
		}
	}
	if key == "" {
		return nil, fmt.Errorf("no host matched with url %v", url)
	}
	properties = allPropertiesFile[key]
	return
}
