package hugo

import (
  "github.com/peter-mount/documentation/tools/util"
  "os"
  "path/filepath"
  "strings"
  "time"
)

type BookHandler func(*Book) error

func WithBook() BookHandler {
  return func(_ *Book) error {
    return nil
  }
}

func (a BookHandler) Then(b BookHandler) BookHandler {
  return func(book *Book) error {
    err := a(book)
    if err != nil {
      return err
    }
    return b(book)
  }
}

func (a BookHandler) Do(book *Book) error {
  return a(book)
}

func (a BookHandler) IfExcelPresent(h func(book *Book, excel util.ExcelBuilder) error) BookHandler {
  return a.Then(func(b *Book) error {
    if b.excel != nil {
      return h(b, b.excel)
    }
    return nil
  })
}

type BookGeneratorHandler func(*Book, string) error

func WithBookGenerator() BookGeneratorHandler {
  return func(_ *Book, _ string) error {
    return nil
  }
}

func (a BookGeneratorHandler) Then(b BookGeneratorHandler) BookGeneratorHandler {
  return func(book *Book, n string) error {
    err := a(book, n)
    if err != nil {
      return err
    }
    return b(book, n)
  }
}

func (a BookHandler) ForEachGenerator(f BookGeneratorHandler) BookHandler {
  return a.Then(func(book *Book) error {
    return book.Generate.ForEach(func(s string) error {
      return f(book, s)
    })
  })
}

type Books []*Book

func (bs Books) ForEach(f BookHandler) error {
  for _, b := range bs {
    err := f(b)
    if err != nil {
      return err
    }
  }
  return nil
}

// Book defines a book that's rendered as pdf
type Book struct {
  ID        string            `yaml:"id"`        // ID of the book, e.g. "bbc" or "6502"
  Title     string            `yaml:"title"`     // Title of book, default title from main page
  Author    string            `yaml:"author"`    // Author of book, default ""
  Copyright string            `yaml:"copyright"` // Copyright
  PDF       *PDF              `yaml:"pdf"`       // Custom PDF config for just this book
  Generate  util.StringSlice  `yaml:"generate"`  // List of generators to run on this book
  modified  time.Time         `yaml:"-"`         // Last Modified time
  excel     util.ExcelBuilder `yaml:"-"`         // Excel builder if present
}

func (b *Book) ContentPath() string {
  return "content/" + b.ID + "/"
}

func (b *Book) Modified() time.Time {
  if b.modified.IsZero() {
    _ = filepath.Walk(b.ContentPath(), func(_ string, info os.FileInfo, _ error) error {
      if !info.IsDir() && info.ModTime().After(b.modified) {
        b.modified = info.ModTime()
      }
      return nil
    })
  }
  return b.modified
}

func replace(s, f string, h func() string) string {
  if strings.Contains(s, f) {
    return strings.ReplaceAll(s, f, h())
  }
  return s
}

func (b *Book) Expand(s string) string {
  // Modified date of content
  s = replace(s, "${modified}", func() string {
    return b.Modified().Format(time.RFC1123)
  })

  // Expand book title or use default from chrome
  s = replace(s, "${title}", func() string {
    t := b.Title
    if t == "" {
      t = "<span class='title'></span>"
    }
    return strings.ReplaceAll(s, "${title}", t)
  })

  s = replace(s, "${author}", func() string {
    return b.Author
  })

  s = replace(s, "${copyright}", func() string {
    return b.Copyright
  })

  return s
}

func (b *Book) GetExcel() util.ExcelBuilder {
  if b.excel == nil {
    b.excel = util.NewExcelBuilder()
  }
  return b.excel
}

func (b *Book) SetExcel(eb util.ExcelBuilder) {
  b.excel = eb
}

// Do runs a function against this instance. When it exits it removes any resources the Book has used freeing up memory.
func (b *Book) Do(f func(*Book) error) error {
  defer func() {
    b.excel = nil
  }()
  return f(b)
}
