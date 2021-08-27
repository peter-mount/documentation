package main

import (
  "github.com/peter-mount/documentation/tools/bbc"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/m6502"
  "github.com/peter-mount/documentation/tools/pdf"
  "github.com/peter-mount/go-kernel"
  "log"
)

func main() {
  err := kernel.Launch(
    // Core modules
    &hugo.Hugo{},
    &pdf.PDF{},
    // The various page/file generators
    &bbc.BBC{},
    &m6502.M6502{},
  )
  if err != nil {
    log.Fatal(err)
  }
}
