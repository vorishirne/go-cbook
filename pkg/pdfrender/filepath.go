package pdfrender

import (
	"fmt"
	"math"
	"net/url"
	"path"
	"strings"
)

// GetRawFilePath takes the Url to be rendered and extracts the final rawPath
// either from url path or fragment
// returns the dir in which the pdf be created and the name of the pdf file.
func (m *Mod) GetRawFilePath(webPageUrl string) (filePath string, err error) {

	// parse base url from mod & the url to render
	webPageUrlObj, err := url.Parse(webPageUrl)
	if err != nil {
		return
	}
	err = m.SetURLProperties(webPageUrl)
	if err != nil {
		return
	}
	baseUrlPath := strings.TrimSpace(strings.Trim(m.State.CurrentURLProperties.Path, "/"))
	if err != nil {
		return
	}

	// get their trimmed paths
	webPageUrlPath := strings.TrimSpace(strings.Trim(webPageUrlObj.Path, "/"))

	// if it has some piece of fragment, then pick fragment as filePath
	// this trick is used to set filePath for blogs, that have no standard filePath structure
	if webPageUrlObj.Fragment != "" {
		filePath = strings.Trim(webPageUrlObj.Fragment, "/")
	} else if strings.HasPrefix(webPageUrlPath, baseUrlPath) {
		// else originally remove basePath from urlPath to get filePath
		filePath = strings.Trim(strings.TrimPrefix(webPageUrlPath, baseUrlPath), "/")
	} else {
		// no options left to get filePath
		err = fmt.Errorf("mod base url %v not valid for %v",
			webPageUrlPath, baseUrlPath)
		return
	}
	// break filePath into dirPath and fileName
	return
}

// this is the super complex and most tough function that had me waited for it
// for days to proceed in this repo

// GetIndexedDir recursively iterate through each subDir in the dirPath
// to return the indexedPath for the respective raw input path, along with a new
// index to be given to a file to be created in that indexedDir.
// It takes care for maintaining a cache of it's results,
//to save calls to subDir that
// have already been visited and also peeks from that cache first,
// before doing iteration on it.
func (m *Mod) GetIndexedDir(rawDir string) (indexedCurrentDir, indexForNewItem string, err error) {

	if visitedDir, ok := m.dirVisited[rawDir]; ok {
		visitedDir.ItemCount++
		indexedCurrentDir = visitedDir.IndexedDirPath
		indexForNewItem = visitedDir.DirIndex + fmt.Sprintf("%02d", visitedDir.ItemCount) + "-"
		return
	}

	slashLastIndex := strings.LastIndexAny(rawDir, "/")
	rawCurrentDir := rawDir[slashLastIndex+1:]
	if slashLastIndex == -1 {
		slashLastIndex = 0
	}
	rawParentDir := rawDir[:slashLastIndex]

	var indexedParentDir, indexForCurrentDir string
	indexedParentDir, indexForCurrentDir, err = m.GetIndexedDir(rawParentDir)
	if err != nil {
		return
	}
	indexedCurrentDir = path.Join(indexedParentDir, indexForCurrentDir+rawCurrentDir)

	indexForNewItem = indexForCurrentDir + fmt.Sprintf("%02d", 1) + "-"
	m.dirVisited[rawDir] = &DirVisited{
		indexedCurrentDir, 1, indexForCurrentDir}
	return
}

func (m *Mod) GetIndexedFilePath(rawPath string) (indexedCurrentFilePath string, err error) {
	if visitedDir, ok := m.dirVisited[rawPath]; ok {
		indexedCurrentFilePath = visitedDir.IndexedDirPath
		return
	}
	indexedCurrentFilePath, _, err = m.GetIndexedDir(rawPath)
	return
}

func (m *Mod) SetURLProperties(url string) (err error) {
	lowestIndex := math.MaxInt64
	key := ""
	for k := range m.State.AllPropertiesFile {
		if index := strings.Index(url, k); index > -1 {
			if index < lowestIndex {
				key = k
				lowestIndex = index
			}
		}
	}
	if key == "" {
		return fmt.Errorf("no host matched with url %v", url)
	}
	m.State.CurrentURLProperties = m.State.AllPropertiesFile[key]
	return
}
