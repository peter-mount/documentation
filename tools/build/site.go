package build

import (
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/jenkinsfile"
	"github.com/peter-mount/go-build/util/makefile"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
)

type Site struct {
	Build     *core.Build     `kernel:"inject"`
	BookShelf *hugo.BookShelf `kernel:"inject"` // Bookshelf
}

func (s *Site) Start() error {
	s.Build.Makefile(50, s.documentation)
	s.Build.Jenkins(50, s.jenkins)
	return nil
}

const (
	siteDir = "public"
)

func (s *Site) documentation(root makefile.Builder, target target.Builder, meta *meta.Meta) {
	genSite := s.Build.Tool("gensite")

	root.Rule(siteDir).
		Mkdir(siteDir)

	root.Rule("site", genSite, siteDir).
		Echo("GEN SITE", siteDir).
		Line(genSite + " -v")
}

func (s *Site) jenkins(builder, node jenkinsfile.Builder) {
	node.Stage("Site").
		Sh("make -f Makefile.gen site")
}
