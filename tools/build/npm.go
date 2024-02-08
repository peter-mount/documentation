package build

import (
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/arch"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
)

type NPM struct {
	Build *core.Build `kernel:"inject"`
}

func (s *NPM) Start() error {
	s.Build.AddExtension(s.extension)
	return nil
}

const (
	npmDir = "node_modules"
)

func (s *NPM) extension(_ arch.Arch, target target.Builder, meta *meta.Meta) {
	if t := target.GetNamedTarget(npmDir); t != nil {
		target.Link(t)
	} else {
		rule := target.Target(npmDir)
		for _, l := range []string{"install postcss postcss-cli", "install autoprefixer", "install"} {
			rule = rule.Echo("NPM", "%s", l).
				Line("npm %s", l)
		}
	}
}
