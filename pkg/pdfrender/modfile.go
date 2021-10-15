package pdfrender

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/watergist/file-engine/reader"
	"os"
	"path"
	"strings"
)

var webPageExtensions = []string{"html", "md", "htm", "hbs", "cms", "php"}

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
	// a book name when compiling books
	BookName string
	// this file is merged with global css rules and passed to wkhtmltopdf
	CssOverrideFile string
	// the base url to be taken away as prefix from the urls' paths
	BaseUrl string
	// base parent dir to all the dirs
	BaseDir string
	// a rune, to mark the current running set of files & dirs, that will be
	// differentiated from the past available files &dirs
	HistPointer string
	// these options are passed to the wkhtmltopdf object
	// could be used for overriding default ones in this repository
	ObjectOptions *json.RawMessage
	// these options are passed to the wkhtmltopdf converter
	// could be used for overriding default ones in this repository
	ConverterOptions *json.RawMessage
	// IsMD is set when the url ends with .md
	// it is set iteratively everytime a url is processed
	IsMD                 bool
	IndexedBookmarkNames bool
	GenBook              bool
	//		   map[rawFilePath]dataForIndex
	dirVisited map[string]*DirVisited
}

// GetMod prepares the Mod struct from mod json file provided
func GetMod(modFilePath string) (m *Mod, err error) {
	modFile, err := os.ReadFile(modFilePath)
	if err != nil {
		return
	}
	m = &Mod{}

	err = json.Unmarshal(modFile, m)
	if err != nil {
		return
	}
	if m.HistPointer == "" {
		m.HistPointer = "0"
	}
	m.BaseDir = path.Join(m.BaseDir, m.HistPointer)
	m.BaseDir = strings.TrimSuffix(m.BaseDir, "/")
	m.BaseDir = strings.TrimSpace(m.BaseDir)

	err = m.GetDirVisited()
	if err != nil {
		return
	}

	return
}

// GetDirVisited tries a lookup for locally stored .gob file that holds the
// index for already existing pdf files.
func (m *Mod) GetDirVisited() (err error) {
	dirIndexGobPath := strings.ReplaceAll(m.BaseDir, "/", "-")
	dirIndexGobPath = strings.Trim(dirIndexGobPath, "-")
	dirIndexGobPath = path.Join(m.BaseDir, dirIndexGobPath) + ".gob"
	_, err = os.Stat(dirIndexGobPath)
	if errors.Is(err, os.ErrNotExist) {
		err = nil

		m.dirVisited = map[string]*DirVisited{"": {
			IndexedDirPath: m.BaseDir,
			ItemCount:      0,
		}}
		return
	}
	gobFile, err := os.Open(dirIndexGobPath)
	defer gobFile.Close()
	if err != nil {
		return
	}
	gobDecoder := gob.NewDecoder(gobFile)
	err = gobDecoder.Decode(&m.dirVisited)
	if err != nil {
		return
	}
	return
}

// SaveDirVisited does exactly opposite
// it does the saving for the file
func (m *Mod) SaveDirVisited() (err error) {
	dirIndexGobPath := strings.ReplaceAll(m.BaseDir, "/", "-")
	dirIndexGobPath = strings.TrimPrefix(dirIndexGobPath, "-")
	dirIndexGobPath = path.Join(m.BaseDir, dirIndexGobPath) + ".gob"
	gobFile, err := os.Create(dirIndexGobPath)
	if err != nil {
		return
	}
	defer gobFile.Close()
	gobEncoder := gob.NewEncoder(gobFile)
	err = gobEncoder.Encode(m.dirVisited)
	return
}

// GetFilePath wraps away the logic to provide & create
// indexed dir from the url path or url fragment
func (m *Mod) GetFilePath(webUrl string) (filePath string, err error) {
	// 1. Extracts raw filepath from url.
	rawFilePath, err := m.GetRawFilePath(webUrl)
	if err != nil {
		return
	}

	// 2. Converts it to indexed one
	filePath, err = m.GetIndexedFilePath(rawFilePath)
	if err != nil {
		return
	}

	parentDir, fileName, err := reader.GetDirPathAndFileName(filePath, true)
	if errors.Is(err, reader.ErrorNoBasePath) {
		err = nil
	}

	//  & creates that dir, even if its already there
	err = os.MkdirAll(parentDir, 0766)
	if err != nil {
		return
	}

	//finally, the file name can have .html/.md/.htm at the end of it, remove them
	extension := ""
	for _, extensions := range webPageExtensions {
		if strings.HasSuffix(fileName, "."+extensions) {
			extension = "." + extensions
			break
		}
	}
	// 3. Returns the final filePath to be used to create pdf at.
	filePath = strings.TrimSuffix(filePath, extension) + ".pdf"
	if extension == ".md" {
		m.IsMD = true
	} else {
		m.IsMD = false
	}

	return
}
