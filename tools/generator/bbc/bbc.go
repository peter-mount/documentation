package bbc

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util/task"
  "github.com/peter-mount/documentation/tools/util/walk"
  "github.com/peter-mount/go-kernel"
  "log"
)

// BBC generates the reference pages in the BBC book.
// These pages are the indices to the various sections like MOS calls, OSByte & OSWord calls etc.
type BBC struct {
  generator *generator.Generator // Generator
  excel     *generator.Excel     // Excel
  extracted bool                 // True once extract() has run
  api       []*Api               // MOS API calls
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

func (b *BBC) Name() string {
  return "BBC"
}

func (b *BBC) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&generator.Generator{})
  if err != nil {
    return err
  }
  b.generator = service.(*generator.Generator)

  service, err = k.AddService(&generator.Excel{})
  if err != nil {
    return err
  }
  b.excel = service.(*generator.Excel)

  return nil
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
        OtherExists("api", b.extractApi).
        OtherExists("osbyte", b.extractOsbyte).
        OtherExists("osword", b.extractOsword).
        Walk(ctx)).
    Walk(book.ContentPath())
}
