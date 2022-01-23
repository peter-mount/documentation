package bbc

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "sort"
)

type Osword struct {
  Call   int                         // Call 0..255
  params map[interface{}]interface{} // Params
  Hex    string                      // Hex code of OSByte Call
  Title  string                      // Title of OSByte Call
  Exit   FunctionParams              // Exit parameters
  Compat Compatibility               // Machine compatibility
}

func (b *BBC) extractOsword(ctx context.Context, _ *hugo.FrontMatter) error {
  return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util.IfMap(e, func(m map[interface{}]interface{}) error {
      if _, ok := util.DecodeInt(m["int"], 0); ok {
        o := &Osword{
          Call:   util.IfMapEntryInt(m, "int"),
          params: m,
          Hex:    util.IfMapEntryString(m, "hex"),
          Title:  util.IfMapEntryString(m, "title"),
        }

        if err := util.IfMapEntry(m, "exit", o.Exit.decode); err != nil {
          return err
        }

        if err := util.IfMapEntry(m, "compatibility", o.Compat.decode); err != nil {
          return err
        }

        b.osword = append(b.osword, o)
      }
      return nil
    })
  })
}

func (b *BBC) writeOswordIndex(ctx context.Context) error {
  book := generator.GetBook(ctx)

  sort.SliceStable(b.osword, func(i, j int) bool {
    return b.osword[i].Call < b.osword[j].Call
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
    book.Modified(),
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "osword", "_index.html"), book.Modified())
}

func (b *BBC) writeOswordTable(ctx context.Context) error {
  book := generator.GetBook(ctx)

  return util.WithTable().
    AsCSV(book.StaticPath("osword.csv"), book.Modified()).
    AsExcel(b.excel.Get(book.ID, book.Modified())).
      Do(&util.Table{
        Title: "osword",
        Columns: []string{
          "Decimal",
          "Hex",
          "Action",
          "Exit A",
          "Exit X",
          "Exit Y",
          "Exit C",
          "BBC",
          "Master",
          "Electron",
          "Other",
        },
        RowCount: len(b.osword),
        GetRow: func(r int) interface{} {
          return b.osword[r]
        },
        Transform: func(i interface{}) []interface{} {
          o := i.(*Osword)
          return []interface{}{
            o.Call,
            o.Hex,
            o.Title,
            o.Exit.A,
            o.Exit.X,
            o.Exit.Y,
            o.Exit.C,
            o.Compat.BBC,
            o.Compat.Master,
            o.Compat.Electron,
            o.Compat.Other,
          }
        },
      })
}
