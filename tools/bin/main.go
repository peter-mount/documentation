package main

import (
	"github.com/peter-mount/documentation/tools/bbc"
	"github.com/peter-mount/documentation/tools/m6502"
	"github.com/peter-mount/go-kernel"
	"log"
)

func main() {
	err := kernel.Launch(
		&m6502.M6502{},
		&bbc.BBC{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
