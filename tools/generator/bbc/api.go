package bbc

import (
  "context"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "log"
  "sort"
  "strconv"
  "strings"
)

type Api struct {
  call     int            // Call 0..255
  Name     string         `yaml:"name"`
  Addr     string         `yaml:"addr"`
  Indirect string         `yaml:"indirect,omitempty"`
  Title    string         `yaml:"title"`
  Entry    FunctionParams `yaml:"entry"`
  Exit     FunctionParams `yaml:"exit"`
  params   interface{}
}

func (b *BBC) extractApi(ctx context.Context, _ *hugo.FrontMatter) error {
  return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util.IfMap(e, func(m map[interface{}]interface{}) error {
      v := &Api{
        Name:     util.DecodeString(m["name"], ""),
        Addr:     util.DecodeString(m["addr"], ""),
        Indirect: util.DecodeString(m["indirect"], ""),
        Title:    util.DecodeString(m["title"], ""),
        params:   e,
      }

      if i, err := strconv.ParseInt(v.Addr, 16, 64); err != nil {
        log.Printf("Failed to parse addr \"%s\" for %s", v.Addr, v.Name)
        return err
      } else {
        v.call = int(i)
      }

      if err := util.IfMapEntry(m, "entry", v.Entry.decode); err != nil {
        return err
      }

      if err := util.IfMapEntry(m, "exit", v.Exit.decode); err != nil {
        return err
      }

      b.api = append(b.api, v)

      return nil
    })
  })
}

func (b *BBC) writeAPIIndex(book *hugo.Book) error {
  sort.SliceStable(b.api, func(i, j int) bool {
    return b.api[i].call < b.api[j].call
  })

  r := Output{Nometa: true}
  for _, o := range b.api {
    r.Api = append(r.Api, o.params)
  }

  return util.ReferenceFileBuilder(
    "MOS API by address",
    "MOS API by address",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "api", "_index.html"), book.Modified())
}

func (b *BBC) writeAPINameIndex(book *hugo.Book) error {
  sort.SliceStable(b.api, func(i, j int) bool {
    return strings.ToLower(b.api[i].Name) < strings.ToLower(b.api[j].Name)
  })

  r := Output{Nometa: true}
  for _, o := range b.api {
    r.Api = append(r.Api, o.params)
  }

  return util.ReferenceFileBuilder(
    "MOS API by name",
    "MOS API by name",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "apiName", "_index.html"), book.Modified())
}

func (b *BBC) writeAPITable(book *hugo.Book) error {
  // Work with a copy as this could be changed once excel runs
  api := append([]*Api{}, b.api...)
  sort.SliceStable(api, func(i, j int) bool {
    return api[i].call < api[j].call
  })

  return util.WithTable().
    AsCSV(util.ReferenceFilename(book.ContentPath(), "api", "api.csv"), book.Modified()).
    AsExcel(book).
      Do(&util.Table{
        Title: "api",
        Columns: []string{
          "Function",
          "Address",
          "Vector",
          "Description",
          "Entry A",
          "Entry X",
          "Entry Y",
          "Exit A",
          "Exit X",
          "Exit Y",
          "Exit C",
        },
        RowCount: len(api),
        GetRow: func(r int) interface{} {
          return api[r]
        },
        Transform: func(i interface{}) []interface{} {
          o := i.(*Api)
          return []interface{}{
            o.Name,
            o.Addr,
            o.Indirect,
            o.Title,
            o.Entry.A,
            o.Entry.X,
            o.Entry.Y,
            o.Exit.A,
            o.Exit.X,
            o.Exit.Y,
            o.Exit.C,
          }
        },
      })
}
