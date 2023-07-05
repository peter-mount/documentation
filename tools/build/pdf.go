package build

import (
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/arch"
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
	s.Build.AddExtension(s.extension)
	return nil
}

func (s *PDF) extension(arch arch.Arch, parentTarget target.Builder, meta *meta.Meta) {
	baseDir := arch.BaseDir(*s.Build.Encoder.Dest)

	if t := parentTarget.GetNamedTarget("pdf"); t != nil {
		parentTarget.Link(t)
		return
	}

	pdfTarget := parentTarget.Target(
		"pdf",
		filepath.Join(baseDir, "bin/genpdf"),
	)

	_ = s.BookShelf.
		Books().
		ForEach(func(book *hugo.Book) error {
			bookTarget := filepath.Join(siteDir, "static/book", book.ID+".pdf")
			pdfTarget.Target(bookTarget, siteDir).
				MkDir(filepath.Dir(bookTarget)).
				Echo("GEN PDF", bookTarget).
				Line(strings.Join([]string{
					filepath.Join(baseDir, "bin", "genpdf"),
					// Uncomment for verbosity
					//"-v",
					book.ID,
					bookTarget,
				}, " "))
			return nil
		})
}
