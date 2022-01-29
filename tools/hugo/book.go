package hugo

import (
  strings2 "github.com/peter-mount/go-kernel/util/strings"
  "os"
  "path"
  "path/filepath"
  "strings"
  "time"
)

type BookHandler func(*Book) error

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
  BookCopyright `yaml:",inline"`                         // Copyright of book
  ID            string               `yaml:"id"`         // ID of the book, e.g. "bbc" or "6502"
  FrontImage    BookCopyright        `yaml:"frontImage"` // Copyright of front image
  Generate      strings2.StringSlice `yaml:"generate"`   // List of generators to run on this book
  modified      time.Time            `yaml:"-"`          // Last Modified time
  contentPath   string
  webPath       string
}

type BookCopyright struct {
  Title     string `yaml:"title"`     // Title of book, default title from main page
  SubTitle  string `yaml:"subTitle"`  // SubTitle
  Author    string `yaml:"author"`    // Author of book, default ""
  SubAuthor string `yaml:"subAuthor"` // SubAuthor of book, default ""
  Copyright string `yaml:"copyright"` // Copyright
}

func (b *Book) ContentPath(s ...string) string {
  if len(s) == 0 {
    return b.contentPath
  }
  return path.Join(append([]string{b.contentPath}, s...)...)
}

func (b *Book) WebPath() string {
  return b.webPath
}

func (b *Book) StaticPath(suffix string) string {
  return path.Join("static/static/book", b.ID+"_"+strings.ReplaceAll(suffix, "/", "_"))
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

// Do runs a function against this instance. When it exits it removes any resources the Book has used freeing up memory.
func (b *Book) Do(f func(*Book) error) error {
  return f(b)
}
