package autodoc

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  util2 "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/walk"
  "log"
)

func (s *Autodoc) extract(ctx context.Context) error {
  book := generator.GetBook(ctx)

  // Only run once per Book ID
  if s.extracted.Contains(book.ID) {
    return nil
  }

  // Prevent us running again for this Book
  s.extracted.Add(book.ID)

  log.Printf("Scanning %s for autodocs", book.ID)

  return walk.NewPathWalker().
    IsFile().
    PathNotContain("/reference/").
    PathHasSuffix(".html").
      Then(hugo.FrontMatterActionOf().
        OtherExists("memorymap", s.extractMemoryMap).
        Context("book", book).
        Walk()).
    Walk(book.ContentPath())
}

func (s *Autodoc) extractMemoryMap(ctx context.Context, _ *hugo.FrontMatter) error {
  headers, err := s.getHeaders(ctx)
  if err != nil {
    return err
  }

  // err is never non nil as headers are attached to the book
  book, _ := s.getBook(ctx)
  headers.Add(&Header{Comment: book.Title})

  return util2.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util2.IfMap(e, func(m map[interface{}]interface{}) error {
      headers.Add(&Header{
        Label:   util2.DecodeString(m["name"], ""),
        Value:   "0x" + util2.DecodeString(m["address"], ""), // Valid as address is always in hex
        Comment: util2.DecodeString(m["desc"], ""),
      })
      return nil
    })
  })
}
