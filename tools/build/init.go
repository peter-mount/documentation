package build

import "github.com/peter-mount/go-kernel/v2"

func init() {
	kernel.Register(
		// Removed hugo as now expected to be installed on the build host
		//&Hugo{},
		&NPM{},
		&PDF{},
		&Site{},
		&Targets{},
	)
}
