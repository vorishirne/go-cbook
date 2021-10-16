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

// PDFMergeState shares data across the merged pdf generation process
type PDFMergeState struct {
	// m.BaseDir set here
	BaseDir string
	// the path to put merged file at
	OutPath string
	// if set true, index in bookmark names will be there
	IndexedBookmarkNames bool
	// number of pages covered so far, so that new bookmark be started from this +1
	currentPageCountReached int
}

func CompileToBook(m *pdfrender.Mod) (err error) {
	// generate a file tree
	tree := list.Read(&[]string{m.BaseDir}, -1)
	// a nested tree is what we want
	tree.GenerateNestedTree = true
	// and also we want it to preserve the order in slice,
	// which will not be the case in map
	tree.GenerateNestedTreeFileOrder = true
	err = tree.UpdateFiles()
	if err != nil {
		return
	}

	if len(tree.NestedPathTree) != 1 {
		return fmt.Errorf("expected only one tree from list tree got %v", len(tree.NestedPathTree))
	}

	// compute the outPath
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
		BaseDir:              m.BaseDir,
		OutPath:              mergedBookPath,
		IndexedBookmarkNames: m.IndexedBookmarkNames,
	}

	// this will have the final parent bookmark
	var b *pdfcpu.Bookmark

	b, err = mState.TriggerEveryPDFMergeInOrder(tree.NestedPathTree[m.BaseDir], nil)
	if err != nil {
		return
	}

	b.Title = m.BookName
	// we only know at the end the final amount of pages
	b.PageThru = mState.currentPageCountReached

	// write bookmark to file
	err = api.AddBookmarksFile(mState.OutPath, mState.OutPath,
		[]pdfcpu.Bookmark{*b}, nil)
	return
}
