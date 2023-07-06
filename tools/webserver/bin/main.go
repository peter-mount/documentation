package main

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/webserver"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(
		&webserver.Webserver{Foreground: true},
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
