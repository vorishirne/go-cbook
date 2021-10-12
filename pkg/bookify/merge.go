package bookify

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/velcrine/goreader/pkg/pdfrender"
)

func CompileToBook(m *pdfrender.Mod) (err error) {
	inFiles := []string{"blog/go/g-01rpc/g-01-01protocols/g-01-01-01gob.pdf", "blog/go/g-02go/g-02-01pkg/g-02-01-02embed-discussion.pdf"}
	err = api.MergeAppendFile(inFiles, "out.pdf", nil)
	if err != nil {
		return
	}
	return
}
