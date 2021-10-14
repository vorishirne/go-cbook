package merge

import (
	"errors"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/velcrine/goreader/pkg/pdfrender"
	"github.com/watergist/file-engine/list"
	"os"
	"path"
	"strings"
)

type PDFMergeState struct {
	BaseDir                 string
	OutPath                 string
	parentBookmark          *pdfcpu.Bookmark
	currentBookmark         *pdfcpu.Bookmark
	currentPageCountReached int
}

func CompileToBook(m *pdfrender.Mod) (err error) {
	tree := list.Read(&[]string{m.BaseDir}, -1)
	tree.GenerateNestedTree = true
	tree.GenerateNestedTreeFileOrder = true

	err = tree.UpdateFiles()
	if err != nil {
		return
	}
	if len(tree.NestedPathTree) != 1 {
		return fmt.Errorf("expected only one tree from list tree got %v", len(tree.NestedPathTree))
	}

	mergedBookPath := strings.ReplaceAll(m.BaseDir, "/", "-")
	mergedBookPath = strings.TrimPrefix(mergedBookPath, "-")
	mergedBookPath = path.Join(m.BaseDir, mergedBookPath) + ".pdf"

	// first delete this file, otherwise pdf cpu will merge new into this
	// instead of overriding
	if _, err = os.Stat(mergedBookPath); err == nil {
		err = os.Remove(mergedBookPath)
		if err != nil {
			return
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return
	}

	mState := &PDFMergeState{
		BaseDir: m.BaseDir,
		OutPath: mergedBookPath,
	}

	var b *pdfcpu.Bookmark
	b, err = mState.TriggerEveryPDFMergeInOrder(tree.NestedPathTree[m.BaseDir], nil)
	if err != nil {
		return
	}
	b.Title = m.BookName
	b.PageThru = mState.currentPageCountReached

	err = api.AddBookmarksFile(mState.OutPath, mState.OutPath,
		[]pdfcpu.Bookmark{*b}, nil)
	return
}
