package build

import (
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
	"path/filepath"
	"strings"
)

type PDF struct {
	Build     *core.Build     `kernel:"inject"`
	BookShelf *hugo.BookShelf `kernel:"inject"` // Bookshelf
}

func (s *PDF) Start() error {
	s.Build.Documentation(s.documentation)
	return nil
}

func (s *PDF) documentation(parentTarget target.Builder, meta *meta.Meta) {
	if t := parentTarget.GetNamedTarget("pdf"); t != nil {
		parentTarget.Link(t)
		return
	}

	genPdf := s.Build.Tool("genpdf")

	pdfTarget := parentTarget.Target("pdf", genPdf)

	_ = s.BookShelf.
		Books().
		ForEach(func(book *hugo.Book) error {
			bookTarget := filepath.Join(siteDir, "static/book", book.ID+".pdf")
			pdfTarget.Target(bookTarget, siteDir, genPdf).
				MkDir(filepath.Dir(bookTarget)).
				Echo("GEN PDF", bookTarget).
				Line(strings.Join([]string{
					genPdf,
					// Uncomment for verbosity
					//"-v",
					book.ID,
					bookTarget,
				}, " "))
			return nil
		})
}
