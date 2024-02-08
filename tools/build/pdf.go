package build

import (
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/jenkinsfile"
	"github.com/peter-mount/go-build/util/makefile"
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
	s.Build.Makefile(100, s.documentation)
	s.Build.Jenkins(100, s.jenkins)
	return nil
}

func (s *PDF) documentation(root makefile.Builder, parentTarget target.Builder, _ *meta.Meta) {
	genPdf := s.Build.Tool("genpdf")

	// Common dependencies for all pdf's
	pdfDependencies := []string{siteDir, genPdf}

	var targets []string

	_ = s.BookShelf.
		Books().
		ForEach(func(book *hugo.Book) error {
			bookTarget := filepath.Join(siteDir, "static/book", book.ID+".pdf")
			targets = append(targets, bookTarget)
			root.Rule(bookTarget, pdfDependencies...).
				Mkdir(filepath.Dir(bookTarget)).
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

	// Now create a rule for those targets
	root.Phony("pdf")
	root.Rule("pdf", s.Build.Tool("genpdf"), "site").
		AddDependency(targets...)
}

func (s *PDF) jenkins(builder, node jenkinsfile.Builder) {
	node.Stage("PDF").
		Sh("make -f Makefile.gen pdf").
		ArchiveArtifacts("public/static/book/*.pdf")
}
