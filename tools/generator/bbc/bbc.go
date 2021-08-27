package bbc

import (
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "log"
  "os"
  "path/filepath"
  "strings"
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
    Register("bbcAPIIndex", b.extract, b.writeAPIIndex).
    Register("bbcAPINameIndex", b.extract, b.writeAPINameIndex).
    Register("bbcOsbyteIndex", b.extract, b.writeOsbyteIndex).
    Register("bbcOswordIndex", b.extract, b.writeOswordIndex)

  return nil
}

func (b *BBC) extract(book *hugo.Book) error {
  if b.extracted {
    return nil
  }

  b.osbyte = nil

  log.Println("Scanning BBC API")
  err := filepath.Walk(book.ContentPath(), func(path string, info os.FileInfo, err error) error {
    return b.extractMeta(book, path, info, err)
  })
  if err != nil {
    return err
  }

  b.extracted = true
  return nil
}

func (b *BBC) extractMeta(book *hugo.Book, path string, info os.FileInfo, err error) error {
  if err != nil {
    return err
  }

  if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") && !strings.Contains(path, "reference") {
    //log.Println("Reading", path)
    m, err := util.LoadFrontMatter("", path)
    if err != nil {
      return err
    }

    if api, exists := m["api"]; exists {
      b.extractApi(api)
    }

    if osbyte, exists := m["osbyte"]; exists {
      b.extractOsbyte(osbyte)
    }

    if osword, exists := m["osword"]; exists {
      b.extractOsword(osword)
    }
  }

  return nil
}
