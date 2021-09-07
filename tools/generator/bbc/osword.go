package bbc

import (
  "context"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "sort"
)

type Osword struct {
  call   int                         // Call 0..255
  params map[interface{}]interface{} // Params
}

func (b *BBC) extractOsword(ctx context.Context, _ *hugo.FrontMatter) error {
  return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util.IfMap(e, func(m map[interface{}]interface{}) error {
      if v, ok := util.DecodeInt(m["int"], 0); ok {
        o := &Osword{
          call:   v,
          params: m,
        }
        b.osword = append(b.osword, o)
      }
      return nil
    })
  })
}

func (b *BBC) writeOswordIndex(book *hugo.Book) error {
  sort.SliceStable(b.osword, func(i, j int) bool {
    return b.osword[i].call < b.osword[j].call
  })

  r := Output{Nometa: true}
  for _, o := range b.osword {
    r.Osword = append(r.Osword, o.params)
  }

  return util.ReferenceFileBuilder(
    "OSWord calls",
    "OSWord &FFF1 calls",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "osword", "_index.html"), book.Modified())
}

func (b *BBC) writeOswordTable(book *hugo.Book) error {
  return util.WithTable().
    AsCSV(util.ReferenceFilename(book.ContentPath(), "osword", "osword.html"), book.Modified()).
    AsExcel(book).
      Do(&util.Table{
        Title:    "osword",
        Columns:  nil,
        RowCount: len(b.osword),
        GetRow: func(r int) interface{} {
          return b.osbyte[r]
        },
        Transform: func(i interface{}) []interface{} {
          //o := i.(*Osword)
          return []interface{}{}
        },
      })
}
