package build

import "github.com/peter-mount/go-kernel/v2"

func init() {
	kernel.Register(
		&NPM{},
		&PDF{},
		&Site{},
		&Targets{},
	)
}
