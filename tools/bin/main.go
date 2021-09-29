package main

import (
  "github.com/peter-mount/documentation/tools/generator/bbc"
  "github.com/peter-mount/documentation/tools/generator/chip"
  "github.com/peter-mount/documentation/tools/generator/m6502"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/pdf"
  "github.com/peter-mount/go-kernel"
  "log"
)

func main() {
  err := kernel.Launch(
    // The various page/file generators. Place these before the core modules
    &bbc.BBC{},
    &m6502.M6502{},
    &chip.Chip{},
    // Core modules. Have these after the generators so they pick up the new content
    &hugo.Hugo{},
    &pdf.PDF{},
  )
  if err != nil {
    log.Fatal(err)
  }
}
