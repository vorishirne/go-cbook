package urlproperties

import (
	"encoding/json"
	"fmt"
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

var allProperties map[string]*URLProperties

func Init(propertiesFiles []string) (err error) {
	if allProperties != nil {
		return
	}
	allProperties = make(map[string]*URLProperties)
	err, _ = structure.LoadJsonFile(DefaultPropertiesFile,
		&allProperties)
	if err != nil {
		return
	}
	for _, v := range propertiesFiles {
		err, _ = structure.LoadJsonFile(v,
			&allProperties)
		if err != nil {
			return
		}
	}
	return
}

func SetURLProperties(url string) (properties *URLProperties, err error) {
	if allProperties == nil {
		err = fmt.Errorf("properties uninitialized")
		return
	}
	if err != nil {
		return
	}
	lowestIndex := math.MaxInt64
	key := ""
	for k := range allProperties {
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
	properties = allProperties[key]
	return
}
