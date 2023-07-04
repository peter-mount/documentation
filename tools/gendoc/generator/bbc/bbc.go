package bbc

import (
	"context"
	"github.com/peter-mount/documentation/tools/gendoc/generator"
	"github.com/peter-mount/documentation/tools/gendoc/hugo"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"log"
)

// BBC generates the reference pages in the BBC book.
// These pages are the indices to the various sections like MOS calls, OSByte & OSWord calls etc.
type BBC struct {
	generator *generator.Generator `kernel:"inject"` // Generator
	excel     *generator.Excel     `kernel:"inject"` // Excel
	extracted bool                 // True once extract() has run
	osbyte    []*Osbyte            // OSBYTE calls
	osword    []*Osword            // OSWORD calls
}

// Output is used for generating the index pages front matter
type Output struct {
	Nometa bool          `yaml:"nometa"`
	Api    []interface{} `yaml:"api,omitempty"`
	Osbyte []interface{} `yaml:"osbyte,omitempty"`
	Osword []interface{} `yaml:"osword,omitempty"`
}

func (b *BBC) Start() error {
	b.generator.
		Register("bbcOsbyteIndex",
			task.Of().
				RunOnce(&b.extracted, b.extract).
				Then(b.writeOsbyteIndex).
				Then(b.writeOsbyteTable)).
		Register("bbcOswordIndex",
			task.Of().
				RunOnce(&b.extracted, b.extract).
				Then(b.writeOswordIndex).
				Then(b.writeOswordTable))

	return nil
}

func (b *BBC) extract(ctx context.Context) error {
	book := generator.GetBook(ctx)

	log.Println("Scanning BBC API")

	return walk.NewPathWalker().
		IsFile().
		PathNotContain("/reference/").
		PathHasSuffix(".html").
		Then(hugo.FrontMatterActionOf().
			OtherExists("osbyte", b.extractOsbyte).
			OtherExists("osword", b.extractOsword).
			Walk(ctx)).
		Walk(book.ContentPath())
}
