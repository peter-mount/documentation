package build

import "github.com/peter-mount/go-kernel/v2"

func init() {
	kernel.Register(
		&Hugo{},
		&NPM{},
		&PDF{},
		&Site{},
		&Targets{},
	)
}
