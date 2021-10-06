package bbc

import (
  "context"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "sort"
)

type Osbyte struct {
  Call   int                         // Call 0..255
  params map[interface{}]interface{} // Params
  Hex    string                      // Hex code of OSByte Call
  Title  string                      // Title of OSByte Call
  Entry  FunctionParams              // Entry parameters
  Exit   FunctionParams              // Exit parameters
  Compat Compatibility               // Machine compatibility
}

type FunctionParams struct {
  A string // A register
  X string // X register
  Y string // Y register
  C string // C Carry Flag
}

func (p *FunctionParams) decode(e interface{}) error {
  return util.IfMap(e, func(m map[interface{}]interface{}) error {
    p.A = util.IfMapEntryString(m, "a")
    p.X = util.IfMapEntryString(m, "x")
    p.Y = util.IfMapEntryString(m, "y")
    p.C = util.IfMapEntryString(m, "c")
    return nil
  })
}

type Compatibility struct {
  BBC      bool   `yaml:"bbc,omitempty"`      // Valid on BBC micro (A, B or B+)
  Master   bool   `yaml:"master,omitempty"`   // Valid on BBC Master 128, 512 or Compact
  Electron bool   `yaml:"electron,omitempty"` // Valid on Acorn Electron
  Other    string `yaml:"other,omitempty"`    // For alternate roms/hardware
}

func (c *Compatibility) decode(e interface{}) error {
  return util.IfMap(e, func(m map[interface{}]interface{}) error {
    c.BBC = util.IfMapEntryBool(m, "bbc")
    c.Master = util.IfMapEntryBool(m, "master")
    c.Electron = util.IfMapEntryBool(m, "electron")
    c.Other = util.IfMapEntryString(m, "other")
    return nil
  })
}

func (b *BBC) extractOsbyte(ctx context.Context, _ *hugo.FrontMatter) error {
  return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util.IfMap(e, func(m map[interface{}]interface{}) error {
      if v, ok := util.DecodeInt(m["int"], 0); ok {
        o := &Osbyte{
          Call:   v,
          params: m,
          Hex:    util.IfMapEntryString(m, "hex"),
          Title:  util.IfMapEntryString(m, "title"),
        }

        err := util.IfMapEntry(m, "entry", o.Entry.decode)
        if err != nil {
          return err
        }

        err = util.IfMapEntry(m, "exit", o.Exit.decode)
        if err != nil {
          return err
        }

        err = util.IfMapEntry(m, "compatibility", o.Compat.decode)
        if err != nil {
          return err
        }

        if o.Entry.A == "" {
          // This is always the case
          o.Entry.A = "&" + o.Hex
        }

        b.osbyte = append(b.osbyte, o)
      }
      return nil
    })
  })
}

func (b *BBC) writeOsbyteIndex(book *hugo.Book) error {
  sort.SliceStable(b.osbyte, func(i, j int) bool {
    return b.osbyte[i].Call < b.osbyte[j].Call
  })

  r := Output{Nometa: true}
  for _, o := range b.osbyte {
    r.Osbyte = append(r.Osbyte, o.params)
  }

  return util.ReferenceFileBuilder(
    "OSByte calls",
    "OSByte &FFF4 calls",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "osbyte", "_index.html"), book.Modified())
}

func (b *BBC) writeOsbyteTable(book *hugo.Book) error {
  return util.WithTable().
    AsCSV(book.StaticPath("osbyte.csv"), book.Modified()).
    AsExcel(b.excel.Get(book.ID, book.Modified())).
      Do(&util.Table{
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
          "Other",
        },
        RowCount: len(b.osbyte),
        GetRow: func(r int) interface{} {
          return b.osbyte[r]
        },
        Transform: func(i interface{}) []interface{} {
          o := i.(*Osbyte)
          return []interface{}{
            o.Call,
            o.Hex,
            o.Title,
            o.Entry.A,
            o.Entry.X,
            o.Entry.Y,
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
