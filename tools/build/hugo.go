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
type Hugo struct {
	Encoder      *core.Encoder `kernel:"inject"`
	Build        *core.Build   `kernel:"inject"`
	HugoRetrieve *bool         `kernel:"flag,retrieve-hugo,retrieve hugo"`
	HugoInstall  *string       `kernel:"flag,install-hugo,install hugo"`
}

func (s *Hugo) Start() error {
	if *s.HugoRetrieve {
		return s.retrieveHugo()
	}

	if *s.HugoInstall != "" {
		return s.installHugo()
	}

	s.Build.AddExtension(s.extension)
	return nil
}

const (
	hugoRepo    = "https://github.com/gohugoio/hugo"
	hugoDir     = "tmp/hugo"
	hugoVersion = "0.108.0"
	hugoSrcUrl  = "https://github.com/gohugoio/hugo/releases/download/v%s/hugo_extended_%s_Linux-64bit.tar.gz"
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
			//Echo("Download", url).
			BuildTool("-retrieve-hugo", "-d", hugoDir)
	}

	dir := filepath.Join(arch.BaseDir(*s.Encoder.Dest), "bin")
	targetFile := filepath.Join(dir, "hugo")
	relTarget, err := filepath.Rel(hugoDir, targetFile)
	if err != nil {
		panic(err)
	}

	rule := target.
		Target(targetFile, hugoDir).
		MkDir(dir).
		Echo("GO BUILD", targetFile).
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
}

// retrieve hugo binaries
// TODO this only works for linux 64bit
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
	util.Label("GIT", "CHECKOUT %s", hugoRepo)
	return util.RunCommand("git", "clone", hugoRepo, *s.Encoder.Dest)
}
