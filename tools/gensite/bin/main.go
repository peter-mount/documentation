package main

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/gensite/generator/autodoc"
	"github.com/peter-mount/documentation/tools/gensite/generator/bbc"
	"github.com/peter-mount/documentation/tools/gensite/generator/chip"
	"github.com/peter-mount/documentation/tools/gensite/generator/m6502"
	"github.com/peter-mount/documentation/tools/gensite/generator/m68k"
	"github.com/peter-mount/documentation/tools/gensite/generator/svg"
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/documentation/tools/gensite/telstar"
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func main() {
	if err := kernel.Launch(
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
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
