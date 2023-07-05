package build

import (
	"github.com/peter-mount/documentation/tools/gendoc/hugo"
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

func (s *PDF) extension(arch arch.Arch, target target.Builder, meta *meta.Meta) {
	baseDir := arch.BaseDir(*s.Build.Encoder.Dest)

	_ = s.BookShelf.
		Books().
		ForEach(func(book *hugo.Book) error {
			bookTarget := filepath.Join(siteDir, "static/book", book.ID+".pdf")
			t := target.GetNamedTarget(bookTarget)
			if t != nil {
				target.Link(t)
			} else {
				target.Target(bookTarget, siteDir).
					MkDir(filepath.Dir(bookTarget)).
					Echo("GEN PDF", bookTarget).
					Line(strings.Join([]string{
						filepath.Join(baseDir, "bin", "genpdf"),
						// Uncomment for verbosity
						//"-v",
						book.ID,
						bookTarget,
					}, " "))
			}
			return nil
		})
}
