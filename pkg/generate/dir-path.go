package generate

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

func (m *Mod) GetRawFilePath(webPageUrl string) (
	fileParentDir, fileName string, err error) {

	webPageUrlObj, err := url.Parse(webPageUrl)
	if err != nil {
		return "", "", err
	}

	baseUrlObj, err := url.Parse(m.BaseUrl)
	if err != nil {
		return "", "", err
	}
	webPageUrlPath := strings.Trim(webPageUrlObj.Path, "/")
	baseUrlPath := strings.Trim(baseUrlObj.Path, "/")
	if webPageUrlObj.Fragment != "" {
		fileParentDir = webPageUrlObj.Fragment
	} else if strings.HasPrefix(webPageUrlPath, baseUrlPath) {
		fileParentDir = strings.TrimPrefix(webPageUrlPath, baseUrlPath)
	} else {
		err = fmt.Errorf("mod base url %v not valid for %v",
			webPageUrlPath, baseUrlPath)
		return
	}
	fileParentDir = strings.Trim(fileParentDir, "/")
	if fileParentDir == "" {
		err = fmt.Errorf("filepath is empty")
		return
	}
	lastSlash := strings.LastIndexAny(fileParentDir, "/")

	fileName = fileParentDir[lastSlash+1:]
	if lastSlash == -1 {
		lastSlash++
	}
	fileParentDir = fileParentDir[:lastSlash]
	if fileParentDir == "" {
		err = fmt.Errorf("no directory to write file to")
	}
	return
}

// index is like 01-02-01abx.pdf

func (m *Mod) GetIndexedDir(rawDir string) (indexedCurrentDir, indexForNewItem string, err error) {

	if visitedDir, ok := m.dirVisited[rawDir]; ok {
		visitedDir.ItemCount++
		indexedCurrentDir = visitedDir.IndexedDirPath
		indexForNewItem = visitedDir.DirIndex + "-" + fmt.Sprintf("%02d", visitedDir.ItemCount)
		return
	}

	slashLastIndex := strings.LastIndexAny(rawDir, "/")
	rawCurrentDir := rawDir[slashLastIndex+1:]
	var indexForCurrentDir string
	if slashLastIndex == -1 {
		m.baseDirLastIndex++
		indexForCurrentDir = m.HistPointer + fmt.Sprintf("%02d", m.baseDirLastIndex)
		indexedCurrentDir = path.Join(m.BaseDir, indexForCurrentDir+rawCurrentDir)
	} else {
		rawParentDir := rawDir[:slashLastIndex]
		var indexedParentDir string
		indexedParentDir, indexForCurrentDir, err = m.GetIndexedDir(rawParentDir)
		if err != nil {
			return
		}
		indexedCurrentDir = path.Join(indexedParentDir, indexForCurrentDir+rawCurrentDir)
	}

	indexForNewItem = indexForCurrentDir + "-" + fmt.Sprintf("%02d", 1)
	m.dirVisited[rawDir] = &DirVisited{
		indexedCurrentDir, 1, indexForCurrentDir}
	return
}
