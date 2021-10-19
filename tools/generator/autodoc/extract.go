package autodoc

import (
  "context"
  "fmt"
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
        OtherExists("api", s.extractApi).
        OtherExists("memorymap", s.extractMemoryMap).
        Walk(ctx)).
    Walk(book.ContentPath())
}

func (s *Autodoc) extractMemoryMap(ctx context.Context, fm *hugo.FrontMatter) error {

  headers := s.getHeaders(ctx)
  if d, exists := fm.Other["description"]; exists {
    _ = headers.Add(&Header{Comment: util2.DecodeString(d, "")})
    log.Println(d)
  }

  return util2.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util2.IfMap(e, func(m map[interface{}]interface{}) error {
      return headers.Add(&Header{
        Label:   util2.DecodeString(m["name"], ""),
        Value:   "0x" + util2.DecodeString(m["address"], ""), // Valid as address is always in hex
        Comment: util2.DecodeString(m["desc"], ""),
      })
    })
  })
}

func (s *Autodoc) extractApi(ctx context.Context, fm *hugo.FrontMatter) error {

  api := s.getApi(ctx)
  // TODO add comment here

  return util2.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util2.IfMap(e, func(m map[interface{}]interface{}) error {
      v := &ApiEntry{
        Name:     util2.DecodeString(m["name"], ""),
        Addr:     "0x" + util2.DecodeString(m["addr"], ""),
        Indirect: util2.DecodeString(m["indirect"], ""),
        Title:    util2.DecodeString(m["title"], ""),
        params:   e,
      }

      if i, ok := util2.DecodeInt(v.Addr, 0); !ok {
        return fmt.Errorf("failed to parse addr %q for %s", v.Addr, v.Name)
      } else {
        v.call = i
      }

      if err := util2.IfMapEntry(m, "entry", v.Entry.decode); err != nil {
        return err
      }

      if err := util2.IfMapEntry(m, "exit", v.Exit.decode); err != nil {
        return err
      }

      api.Add(v)

      return nil
    })
  })
}
