package merge

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/watergist/file-engine/list"
	"github.com/watergist/file-engine/reader"
	"regexp"
	"strings"
)

func (mState *PDFMergeState) TriggerEveryPDFMergeInOrder(nTree *list.NestedPathTree, parentBookmark *pdfcpu.Bookmark) (bookmark *pdfcpu.Bookmark, err error) {
	bookmark = &pdfcpu.Bookmark{
		PageFrom: mState.currentPageCountReached + 1,
		Bold:     true,
		Italic:   true,
		Parent:   parentBookmark,
	}
	if len(nTree.FileOrder) == 0 {
		return nil, nil
	}
	for _, pdfFileOrDir := range nTree.FileOrder {
		var (
			lastItem string
			b        *pdfcpu.Bookmark
		)
		_, lastItem, err = reader.GetDirPathAndFileName(pdfFileOrDir.Path, true)
		if lastItem == "" {
			return
		} else if !mState.IndexedBookmarkNames {
			reg := regexp.MustCompile(`\w(?:\d\d-)*(.*)\.pdf`)
			lastItem = reg.ReplaceAllString(lastItem, "${1}")
		}
		if pdfFileOrDir.NestedPathTree == nil {
			b, err = mState.MergePdf(pdfFileOrDir.Path, bookmark)
			if err != nil {
				return
			}
		} else {
			b, err = mState.TriggerEveryPDFMergeInOrder(pdfFileOrDir.NestedPathTree, bookmark)
			if err != nil {
				return
			}
		}
		if b != nil {
			b.Title = lastItem
			bookmark.Children = append(bookmark.Children, *b)
		}
	}
	bookmark.PageThru = mState.currentPageCountReached
	return
}

func (mState *PDFMergeState) MergePdf(filePath string, parentBookmark *pdfcpu.Bookmark) (bookmark *pdfcpu.Bookmark, err error) {
	if !strings.HasSuffix(filePath, ".pdf") || filePath == mState.OutPath {
		return
	}

	bookmark = &pdfcpu.Bookmark{
		PageFrom: mState.currentPageCountReached + 1,
		Bold:     true,
		Italic:   true,
		Parent:   parentBookmark,
	}
	nPages, err := api.PageCountFile(filePath)
	if err != nil {
		return
	}
	mState.currentPageCountReached += nPages
	bookmark.PageThru = mState.currentPageCountReached
	err = api.MergeAppendFile([]string{filePath}, mState.OutPath, nil)
	return
}
