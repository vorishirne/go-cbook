package convert

import (
	"github.com/velcrine/eink-pages/pkg/generate"
	"strings"
)

func ReadUrlFile(urlFilePath string, modeFilePath string) (err error) {
	m, err := generate.AllPagesMod(modeFilePath)
	if err != nil {
		return
	}
	readLineCallback := func(line string) (err error) {
		if strings.TrimSpace(line) != "" {
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

	return
}
