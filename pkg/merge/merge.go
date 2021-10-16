package merge

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/watergist/file-engine/list"
	"github.com/watergist/file-engine/reader"
	"regexp"
	"strings"
)

// TriggerEveryPDFMergeInOrder is a recursive function that merges all the pdf
// in current directory and returns the bookmark for that directory,

func (mState *PDFMergeState) TriggerEveryPDFMergeInOrder(nTree *list.NestedPathTree, parentBookmark *pdfcpu.Bookmark) (bookmark *pdfcpu.Bookmark, err error) {
	// new bookmark for this directory
	bookmark = &pdfcpu.Bookmark{
		PageFrom: mState.currentPageCountReached + 1,
		Bold:     true,
		Italic:   true,
		Parent:   parentBookmark,
	}
	// if there is some file/dir in this dir
	if len(nTree.FileOrder) == 0 {
		return nil, nil
	}

	for _, pdfFileOrDir := range nTree.FileOrder {
		var (
			// the name of (pdf file) or (last folder) for this path
			lastItem string
			// the bookmark returned by this (pdf file) or folder
			b *pdfcpu.Bookmark
		)

		_, lastItem, err = reader.GetDirPathAndFileName(pdfFileOrDir.Path, true)
		if lastItem == "" {
			return
			// if index not need to be included, regex it out
		} else if !mState.IndexedBookmarkNames {
			reg := regexp.MustCompile(`\w(?:\d\d-)*(.*)`)
			lastItem = reg.ReplaceAllString(lastItem, "${1}")
			if lastItem == "" {
				err = fmt.Errorf("%v is not a valid pattern, for indexed", pdfFileOrDir.Path)
				return
			}
			lastItem = strings.TrimSuffix(lastItem, ".pdf")
		}
		// check if current path is for a file
		if pdfFileOrDir.NestedPathTree == nil {
			// merge this pdf to output path and also get the bookmark for it to add in parent
			b, err = mState.MergePdf(pdfFileOrDir.Path, bookmark)
			if err != nil {
				return
			}
		} else {
			// merge this entire dir to output path and also get the bookmark for it to add in parent
			b, err = mState.TriggerEveryPDFMergeInOrder(pdfFileOrDir.NestedPathTree, bookmark)
			if err != nil {
				return
			}
		}
		if b != nil {
			b.Title = lastItem
			// update this bookmark in current directory's bookmark
			bookmark.Children = append(bookmark.Children, *b)
		}
	}
	// only now we will know how many pages were there.
	bookmark.PageThru = mState.currentPageCountReached
	return
}

// MergePdf finally merges the pdf "file" in the final output file path and
// returns the bookmark for that file that can even be nil for skipped case
func (mState *PDFMergeState) MergePdf(filePath string, parentBookmark *pdfcpu.Bookmark) (bookmark *pdfcpu.Bookmark, err error) {
	if !strings.HasSuffix(filePath, ".pdf") || filePath == mState.OutPath {
		return
	}

	// new bookmark for pdf file
	bookmark = &pdfcpu.Bookmark{
		PageFrom: mState.currentPageCountReached + 1,
		Bold:     true,
		Italic:   true,
		Parent:   parentBookmark,
	}
	// get number of pages and update global count
	nPages, err := api.PageCountFile(filePath)
	if err != nil {
		return
	}
	mState.currentPageCountReached += nPages
	bookmark.PageThru = mState.currentPageCountReached

	// merge
	err = api.MergeAppendFile([]string{filePath}, mState.OutPath, nil)
	return
}
