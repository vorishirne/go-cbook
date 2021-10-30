package config

import (
	"errors"
	"github.com/vorishirne/go-cbook/pkg/merge"
	"github.com/vorishirne/go-cbook/pkg/pdfrender"
	"github.com/vorishirne/go-cbook/pkg/render"
	"github.com/vorishirne/goreader/reader"
	"github.com/vorishirne/goreader/reader/structure"
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
	if err != nil {
		return
	}
	structure.WriteYaml(
		path.Join(m.BaseDir, m.HistPointer+"err.yaml"), errMap)
	err = m.SaveDirVisited()
	// so errMap is a dict for key: url, v: error returned while rendering it.
	if err != nil {
		return
	}
	if !m.DisableBookGen {
		err = merge.CompileToBook(m)
	}
	return
}
