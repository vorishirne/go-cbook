package convert

import (
	"github.com/velcrine/eink-pages/pkg/generate"
	"github.com/watergist/file-engine/reader"
	"github.com/watergist/file-engine/reader/structures"
	"path"
	"strings"
)

func ReadUrlFile(urlFilePath string, modeFilePath string) (err error) {
	m, err := generate.AllPagesMod(modeFilePath)
	if err != nil {
		return
	}
	readLineCallback := func(line string) (err error) {
		line = strings.TrimSpace(line)
		if line != "" {
			var path string
			path, err = m.GetFilePath(line)
			if err != nil {
				return
			}
			err = m.GenPDF(line, path)
			if err != nil {
				return
			}
		}
		return
	}
	errMap, err := reader.CallbackOnEachLine(urlFilePath, readLineCallback)
	structures.WriteYaml(
		path.Join(m.BaseDir, m.HistPointer+"err.yaml"), errMap)
	return
}
