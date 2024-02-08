package build

import (
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
)

type Site struct {
	Build     *core.Build     `kernel:"inject"`
	BookShelf *hugo.BookShelf `kernel:"inject"` // Bookshelf
}

func (s *Site) Start() error {
	s.Build.Documentation(s.documentation)
	return nil
}

const (
	siteDir = "public"
)

func (s *Site) documentation(target target.Builder, meta *meta.Meta) {
	if t := target.GetNamedTarget(siteDir); t != nil {
		target.Link(t)
		return
	}

	genSite := s.Build.Tool("gensite")

	target.Target(siteDir, genSite).
		Echo("GEN SITE", siteDir).
		Line(genSite)

}
