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

type Site struct {
	Build     *core.Build     `kernel:"inject"`
	BookShelf *hugo.BookShelf `kernel:"inject"` // Bookshelf
}

func (s *Site) Start() error {
	s.Build.AddExtension(s.extension)
	return nil
}

const (
	siteDir = "public"
)

func (s *Site) extension(arch arch.Arch, target target.Builder, meta *meta.Meta) {
	baseDir := arch.BaseDir(*s.Build.Encoder.Dest)

	if t := target.GetNamedTarget(siteDir); t != nil {
		target.Link(t)
	} else {
		target.Target(
			siteDir,
			filepath.Join(arch.BaseDir(*s.Build.Encoder.Dest), "bin", "gensite"),
			//filepath.Join(arch.BaseDir(*s.Build.Encoder.Dest), "bin", "hugo"),
		).
			Echo("GEN SITE", siteDir).
			Line(strings.Join([]string{
				filepath.Join(baseDir, "bin", "gensite"),
				// Uncomment for verbosity
				//"-v"
			}, " "))
	}

}
