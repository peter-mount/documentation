package build

import (
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util"
	"github.com/peter-mount/go-build/util/arch"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
	"os"
	"path/filepath"
)

// Hugo extension fetches the version of Hugo required
//
// # TODO this only works for linux 64bit
//
// it fails with:
// imports github.com/bep/golibsass/internal/libsass: build constraints exclude all Go files in go/pkg/mod/github.com/bep/golibsass@v1.1.1/internal/libsass
// when cross compiling
type Hugo struct {
	Encoder      *core.Encoder `kernel:"inject"`
	Build        *core.Build   `kernel:"inject"`
	HugoRetrieve *bool         `kernel:"flag,retrieve-hugo,retrieve hugo"`
}

func (s *Hugo) Start() error {
	if *s.HugoRetrieve {
		return s.retrieveHugo()
	}

	s.Build.AddExtension(s.extension)
	return nil
}

const (
	hugoRepo = "https://github.com/gohugoio/hugo"
	hugoDir  = "tmp/hugo"
)

func (s *Hugo) extension(arch arch.Arch, target target.Builder, meta *meta.Meta) {
	// We have just one instance of this target so declare it only once.
	// For subsequent calls just link it so the platform gets it as a dependency
	if t := target.GetNamedTarget(hugoDir); t != nil {
		target.Link(t)
	} else {
		target.
			Target(hugoDir).
			MkDir("tmp").
			Echo("GIT", hugoRepo).
			BuildTool("-retrieve-hugo", "-d", hugoDir)
	}

	dir := filepath.Join(arch.BaseDir(*s.Encoder.Dest), "bin")
	targetFile := filepath.Join(dir, "hugo")
	cachedTarget := filepath.Join(arch.BaseDir("tmp"), "hugo")
	relTarget, err := filepath.Rel(hugoDir, cachedTarget)
	if err != nil {
		panic(err)
	}

	rule := target.
		Target(cachedTarget, hugoDir).
		MkDir(dir).
		Echo("GO BUILD", cachedTarget).
		Line(`cd %q;\`, hugoDir)

	if arch.GOARM == "" {
		rule = rule.Line("GCO_ENABLED=0 GOOS=%s GOARCH=%s go build -tags extended -o %q main.go",
			arch.GOOS,
			arch.GOARCH,
			relTarget)
	} else {
		rule = rule.Line("GCO_ENABLED=0 GOOS=%s GOARCH=%s GOARM=%s go build -tags extended -o %q main.go",
			arch.GOOS,
			arch.GOARCH,
			arch.GOARM,
			relTarget)
	}

	target.
		Target(targetFile, cachedTarget).
		MkDir(dir).
		Echo("INSTALL", targetFile).
		BuildTool("-copyfile", cachedTarget, "-d", targetFile)

}

func (s *Hugo) retrieveHugo() error {
	_, err := os.Stat(*s.Encoder.Dest)
	switch {
	case err == nil:
		return nil
	case os.IsNotExist(err):
		return s.checkout()
	default:
		return err
	}
}

func (s *Hugo) checkout() error {
	util.Label("GIT", "CLONE %s", hugoRepo)
	return util.RunCommand("git", "clone", hugoRepo, *s.Encoder.Dest)
}
