package convert

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

func (m *Mod) GetFilePath(webUrl string) (dirPath string, err error) {
	rawDirPath, fileName, err := m.GetRawFilePath(webUrl)
	if err != nil {
		return
	}
	parentDir, itemIndex, err := m.EnsureNumberedDir(rawDirPath)
	if err != nil {
		return
	}
	dirPath = path.Join(parentDir, itemIndex+fileName+".pdf")
	return
}

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

func (m *Mod) EnsureNumberedDir(rawDirPath string) (parentDir, fileCounter string, err error) {

	if info, ok := m.dirVisited[rawDirPath]; ok {
		info.CountReached++
		parentDir = info.DirPath
		fileCounter = info.LastCounterString + "-" + fmt.Sprintf("%03d", info.CountReached)
		return
	}

	lastSlash := strings.LastIndexAny(rawDirPath, "/")
	if lastSlash == -1 {
		m.BaseDirCountReached++
		newFileCounter := fmt.Sprintf("%03d", m.BaseDirCountReached)
		parentDir = path.Join(m.BaseDir, newFileCounter+"-"+rawDirPath)

		//common code
		err = os.MkdirAll(parentDir, 0766)
		if err != nil {
			return
		}
		fileCounter = newFileCounter + "-" + fmt.Sprintf("%03d", 1)
		m.dirVisited[rawDirPath] = &DirVisited{
			parentDir, 1, newFileCounter}
	} else {
		suffixDirPath := rawDirPath[lastSlash+1:]
		prefixDirPath := rawDirPath[:lastSlash]
		var newDirName, newFileCounter string
		newDirName, newFileCounter, err = m.EnsureNumberedDir(prefixDirPath)
		if err != nil {
			return
		}
		parentDir = path.Join(newDirName, newFileCounter+"-"+suffixDirPath)

		// common code
		err = os.Mkdir(parentDir, 0766)
		if err != nil {
			return
		}
		fileCounter = newFileCounter + "-" + fmt.Sprintf("%03d", 1)
		m.dirVisited[rawDirPath] = &DirVisited{
			parentDir, 1, newFileCounter}
	}
	return
}
