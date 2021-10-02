package convert

import (
	"errors"
	"github.com/velcrine/eink-pages/pkg/generate"
	"github.com/watergist/file-engine/reader"
	"github.com/watergist/file-engine/reader/structures"
	"io/fs"
	"os"
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
			var filePath string
			filePath, err = m.GetFilePath(line)
			if err != nil {
				return
			}
			_, e := os.Stat(filePath)
			if errors.Is(e, fs.ErrNotExist) {
				err = m.GenPDF(line, filePath)
			}
		}
		return
	}
	errMap, err := reader.CallbackOnEachLine(urlFilePath, readLineCallback)
	structures.WriteYaml(
		path.Join(m.BaseDir, m.HistPointer+"err.yaml"), errMap)
	return
}
