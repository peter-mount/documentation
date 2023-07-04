package build

import (
	"github.com/peter-mount/go-build/core"
)

type Targets struct {
	Build *core.Build `kernel:"inject"`
}

func (s *Targets) Start() error {
	s.Build.AddCleanDirectory("public")
	return nil
}
