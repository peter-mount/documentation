package main

import (
	"github.com/peter-mount/documentation/tools/gendoc/generator/autodoc"
	"github.com/peter-mount/documentation/tools/gendoc/generator/bbc"
	"github.com/peter-mount/documentation/tools/gendoc/generator/chip"
	"github.com/peter-mount/documentation/tools/gendoc/generator/m6502"
	"github.com/peter-mount/documentation/tools/gendoc/generator/m68k"
	"github.com/peter-mount/documentation/tools/gendoc/generator/svg"
	"github.com/peter-mount/documentation/tools/gendoc/hugo"
	"github.com/peter-mount/documentation/tools/gendoc/telstar"
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
		// Core modules. Have these after the generators, so they pick up the new content
		&hugo.Hugo{},
		// For now remove PDF generation, that will become a standalone tool
		//&pdf.PDF{},
		//&latex.LaTeX{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
