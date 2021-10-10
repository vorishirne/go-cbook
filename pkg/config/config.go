package config

import (
	"errors"
	"github.com/velcrine/eink-pages/pkg/pdfrender"
	"github.com/velcrine/eink-pages/pkg/render"
	"github.com/watergist/file-engine/reader"
	"github.com/watergist/file-engine/reader/structure"
	"io/fs"
	"os"
	"path"
	"strings"
)

// RenderFromUrlFile iterates to every link present in .url file
// and run rendering process on each.
func RenderFromUrlFile(urlFilePath string, modeFilePath string) (err error) {
	m, err := pdfrender.GetMod(modeFilePath)
	if err != nil {
		return
	}

	// this is a callback to provide to reader.CallbackOnEachLine function
	// which runs rendering process for each input line(url) if its non-empty
	// also skips already present pdf.
	readLineCallback := func(line string) (err error) {
		line = strings.TrimSpace(line)
		if line != "" {
			var filePath string
			// the final path and name with which the pdf will be created
			filePath, err = m.GetFilePath(line)
			if err != nil {
				return
			}
			_, e := os.Stat(filePath)
			// if it not exists
			if errors.Is(e, fs.ErrNotExist) {
				err = render.Do(m, line, filePath)
			}
		}
		return
	}
	errMap, err := reader.CallbackOnEachLine(urlFilePath, readLineCallback)
	// so errMap is a dict for key: url, v: error returned while rendering it.
	structure.WriteYaml(
		path.Join(m.BaseDir, m.HistPointer+"err.yaml"), errMap)
	return
}
