package build

import (
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/arch"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
	"path/filepath"
)

type Site struct {
	Build *core.Build `kernel:"inject"`
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
	dir := filepath.Join(baseDir, siteDir)

	target.Target(dir).
		Echo("GENDOC", dir).
		Line(filepath.Join(baseDir, "bin", "gendoc") + " -vt")
}
