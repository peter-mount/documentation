package bbc

import (
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "sort"
)

type Osbyte struct {
  call   int                         // Call 0..255
  params map[interface{}]interface{} // Params
  Hex    string
  Title  string
  Entry  FunctionParams
  Exit   FunctionParams
  Compat Compatibility
}

type FunctionParams struct {
  A string
  X string
  Y string
  C string
}

func (p *FunctionParams) decode(e interface{}) {
  util.IfMap(e, func(m map[interface{}]interface{}) {
    util.IfMapEntry(m, "a", func(i interface{}) {
      p.A = util.DecodeString(i, "")
    })

    util.IfMapEntry(m, "x", func(i interface{}) {
      p.X = util.DecodeString(i, "")
    })

    util.IfMapEntry(m, "y", func(i interface{}) {
      p.Y = util.DecodeString(i, "")
    })

    util.IfMapEntry(m, "c", func(i interface{}) {
      p.C = util.DecodeString(i, "")
    })
  })
}

type Compatibility struct {
  BBC      bool
  Master   bool
  Electron bool
}

func (c *Compatibility) decode(e interface{}) {
  util.IfMap(e, func(m map[interface{}]interface{}) {
    c.BBC = util.IfMapEntryBool(m, "bbc")
    c.Master = util.IfMapEntryBool(m, "master")
    c.Electron = util.IfMapEntryBool(m, "electron")
  })
}

func (b *BBC) extractOsbyte(osbyte interface{}) {
  util.ForEachInterface(osbyte, func(e interface{}) {
    util.IfMap(e, func(m map[interface{}]interface{}) {
      if v, ok := util.DecodeInt(m["int"], 0); ok {
        o := &Osbyte{
          call:   v,
          params: m,
          Hex:    util.IfMapEntryString(m, "hex"),
          Title:  util.IfMapEntryString(m, "title"),
        }

        util.IfMapEntry(m, "entry", o.Entry.decode)
        util.IfMapEntry(m, "exit", o.Exit.decode)
        util.IfMapEntry(m, "compatibility", o.Compat.decode)

        b.osbyte = append(b.osbyte, o)
      }
    })
  })
}

func (b *BBC) writeOsbyteIndex(book *hugo.Book) error {
  sort.SliceStable(b.osbyte, func(i, j int) bool {
    return b.osbyte[i].call < b.osbyte[j].call
  })

  r := Output{Nometa: true}
  for _, o := range b.osbyte {
    r.Osbyte = append(r.Osbyte, o.params)
  }

  err := util.ReferenceFileBuilder(
    "OSByte calls",
    "OSByte &FFF4 calls",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "osbyte", "_index.html"), book.Modified())
  if err != nil {
    return err
  }

  t := util.Table{
    Title: "osbyte",
    Columns: []string{
      "Decimal",
      "Hex",
      "Action",
      "Entry A",
      "Entry X",
      "Entry Y",
      "Exit A",
      "Exit X",
      "Exit Y",
      "Exit C",
      "BBC",
      "Master",
      "Electron",
    },
    RowCount: len(b.osbyte),
    GetRow: func(r int) interface{} {
      return b.osbyte[r]
    },
    Transform: func(i interface{}) []interface{} {
      o := i.(*Osbyte)
      return []interface{}{
        o.call,
        o.Hex,
        o.Title,
        o.call, // Entry.A is always the call for osbyte ;-)
        o.Entry.X,
        o.Entry.Y,
        o.Exit.A,
        o.Exit.X,
        o.Exit.Y,
        o.Exit.C,
        o.Compat.BBC,
        o.Compat.Master,
        o.Compat.Electron,
      }
    },
  }

  return util.NewCSVBuilder().
    Headings(t.Columns...).
    ImportFrom(t).
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "osbyte", "osbyte.csv"), book.Modified())
}
