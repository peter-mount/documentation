package main

import (
  "github.com/peter-mount/documentation/tools/generator/autodoc"
  "github.com/peter-mount/documentation/tools/generator/bbc"
  "github.com/peter-mount/documentation/tools/generator/chip"
  "github.com/peter-mount/documentation/tools/generator/m6502"
  "github.com/peter-mount/documentation/tools/generator/m68k"
  "github.com/peter-mount/documentation/tools/generator/svg"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/pdf"
  "github.com/peter-mount/documentation/tools/telstar"
  "github.com/peter-mount/go-kernel/v2"
  "log"
)

func main() {
  err := kernel.Launch(
    // The various page/file generators. Place these before the core modules
    &bbc.BBC{},
    &m6502.M6502{},
    &m68k.M68k{},
    &chip.Chip{},
    &autodoc.Autodoc{},
    &svg.SVG{},
    &telstar.Service{},
    // Core modules. Have these after the generators so they pick up the new content
    &hugo.Hugo{},
    &pdf.PDF{},
  )
  if err != nil {
    log.Fatal(err)
  }
}
