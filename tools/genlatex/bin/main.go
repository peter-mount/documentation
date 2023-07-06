package main

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(
		&genlatex.LaTeX{},
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
