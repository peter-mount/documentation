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
    // Content generators first
    &m6502.M6502{},
    &bbc.BBC{},
    // Hugo & site generators last
    &hugo.Hugo{},
    &pdf.PDF{},
  )
  if err != nil {
    log.Fatal(err)
  }
}
