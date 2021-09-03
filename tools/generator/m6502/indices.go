package m6502

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "sort"
)

func (s *M6502) writeOpsIndex(book *hugo.Book) error {
  sort.SliceStable(s.opCodes, func(i, j int) bool {
    return s.opCodes[i].Order < s.opCodes[j].Order
  })

  return s.writeFile(
    book,
    "opcodes",
    "Instruction List by opcode",
    "6502 instructions by hex opcode",
  )
}

func (s *M6502) writeOpsHexIndex(book *hugo.Book) error {
  sort.SliceStable(s.opCodes, func(i, j int) bool {
    if s.opCodes[i].Op == s.opCodes[j].Op {
      return s.opCodes[i].Addressing < s.opCodes[j].Addressing
    }
    return s.opCodes[i].Op < s.opCodes[j].Op
  })

  return s.writeFile(
    book,
    "instructions",
    "Instruction List by name",
    "6502 instructions by name",
  )
}

func (s *M6502) writeFile(book *hugo.Book, name, title, desc string) error {
  err := util.ReferenceFileBuilder(
    title,
    desc,
    "manual",
    10,
  ).
      Then(func(a util.StringSlice) (util.StringSlice, error) {
        a = append(a, "codes:")

        for _, op := range s.opCodes {
          a = append(a,
            "  - code: \""+op.Code+"\"",
            "    op: \""+op.Op+"\"",
            "    addressing: "+op.Addressing,
            "    compatibility:",
          )

          for k, _ := range op.Compatibility {
            a = append(a, fmt.Sprintf("      %v: true", k))
          }

          a = op.Bytes.append("    ", "bytes", a)
          a = op.Cycles.append("    ", "cycles", a)
        }

        a = append(a, "notes:")
        for _, n := range s.notes.Notes {
          a = append(a, fmt.Sprintf("  - \"%s\"", n.Value))
        }
        return a, nil
      }).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), name, "_index.html"), book.Modified())
  if err != nil {
    return err
  }

  t := util.Table{
    Title: name,
    Columns: []string{
      "Decimal",
      "Hex",
      "Instruction",
      "Addressing",
      "Len(byte)",
      "Cycles",
    },
    RowCount: len(s.opCodes),
    GetRow: func(r int) interface{} {
      return s.opCodes[r]
    },
    Transform: func(i interface{}) []interface{} {
      o := i.(*Opcode)
      return []interface{}{
        o.Order,
        o.Code,
        o.Op,
        o.Addressing,
        o.Bytes.Value,
        o.Cycles.Value,
      }
    },
  }

  return t.AsCSV().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), name, name+".csv"), book.Modified())
}
