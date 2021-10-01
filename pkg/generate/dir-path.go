package generate

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

func (m *Mod) GetRawFilePath(webUrl string) (dirPath, fileName string, err error) {

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
