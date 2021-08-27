package bbc

import (
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "log"
  "os"
  "path/filepath"
  "strings"
)

type BBC struct {
  generator *hugo.Generator // Generator
  extracted bool            // True once extract() has run
  api       []*Api          // MOS API calls
  osbyte    []*Osbyte       // OSBYTE calls
  osword    []*Osword       // OSWORD calls
}

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
  service, err := k.AddService(&hugo.Generator{})
  if err != nil {
    return err
  }
  b.generator = service.(*hugo.Generator)

  return nil
}

func (b *BBC) Start() error {
  b.generator.
    RegisterCompound("bbcAPIIndex", b.extract, b.writeAPIIndex).
    RegisterCompound("bbcAPINameIndex", b.extract, b.writeAPINameIndex).
    RegisterCompound("bbcOsbyteIndex", b.extract, b.writeOsbyteIndex).
    RegisterCompound("bbcOswordIndex", b.extract, b.writeOswordIndex)

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
