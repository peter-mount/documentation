package bbc

import (
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util/walk"
  "github.com/peter-mount/go-kernel"
  "log"
  "os"
)

// BBC generates the reference pages in the BBC book.
// These pages are the indices to the various sections like MOS calls, OSByte & OSWord calls etc.
type BBC struct {
  generator *generator.Generator // Generator
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

  return nil
}

func (b *BBC) Start() error {
  b.generator.
    Register("bbcAPIIndex", generator.GeneratorHandlerOf(b.extract, b.writeAPIIndex)).
    Register("bbcAPINameIndex", generator.GeneratorHandlerOf(b.extract, b.writeAPINameIndex)).
    Register("bbcOsbyteIndex", generator.GeneratorHandlerOf(b.extract, b.writeOsbyteIndex)).
    Register("bbcOswordIndex", generator.GeneratorHandlerOf(b.extract, b.writeOswordIndex))

  return nil
}

func (b *BBC) extract(book *hugo.Book) error {
  if b.extracted {
    return nil
  }

  log.Println("Scanning BBC API")

  err := walk.NewPathWalker().
    IsFile().
    PathNotContain("/reference/").
    PathHasSuffix(".html").
    Then(b.extractMeta).
    Walk(book.ContentPath())
  if err != nil {
    return err
  }

  b.extracted = true
  return nil
}

func (b *BBC) extractMeta(path string, _ os.FileInfo) error {
  fm := hugo.FrontMatter{}
  err := fm.LoadFrontMatter("", path)
  if err != nil {
    return err
  }

  if api, exists := fm.Other["api"]; exists {
    b.extractApi(api)
  }

  if osbyte, exists := fm.Other["osbyte"]; exists {
    b.extractOsbyte(osbyte)
  }

  if osword, exists := fm.Other["osword"]; exists {
    b.extractOsword(osword)
  }

  return nil
}
