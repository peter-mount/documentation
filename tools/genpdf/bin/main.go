package main

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/genpdf/pdf"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(
		&pdf.PDF{},
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
