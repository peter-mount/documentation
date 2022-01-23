package m6502

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  strings2 "github.com/peter-mount/documentation/tools/util/strings"
  "github.com/peter-mount/go-kernel/util/task"
  "sort"
  "strconv"
  "strings"
)

func delayOpTask(t task.Task) task.Task {
  return func(ctx context.Context) error {
    task.GetQueue(ctx).AddPriorityTask(tools.PriorityIndices,
      task.Of(t).
        WithContext(ctx, generator.BookKey))
    return nil
  }
}

func decodeOpcode(s string) int64 {
  s = strings.ReplaceAll(s, "nn", "00")
  for len(s) < 8 {
    s = s + "00"
  }
  if i, err1 := strconv.ParseInt(s, 16, 64); err1 == nil {
    return i
  }
  return -1
}
func (s *M6502) writeOpsIndex(ctx context.Context) error {
  _ = s.autodoc.GetApi(ctx)

  book := generator.GetBook(ctx)
  inst := s.Instructions(book)

  sort.SliceStable(inst.opCodes, func(i, j int) bool {
    a := inst.opCodes[i].Code
    b := inst.opCodes[j].Code

    return decodeOpcode(a) < decodeOpcode(b)
  })

  return s.writeFile(
    book,
    inst,
    "codes",
    "opcodes",
    "Instruction List by opcode",
    "",
  )
}

func (s *M6502) writeOpsHexIndex(ctx context.Context) error {
  _ = s.autodoc.GetApi(ctx)

  book := generator.GetBook(ctx)
  inst := s.Instructions(book)

  sort.SliceStable(inst.opCodes, func(i, j int) bool {
    if inst.opCodes[i].Op == inst.opCodes[j].Op {
      return inst.opCodes[i].Addressing < inst.opCodes[j].Addressing
    }
    return inst.opCodes[i].Op < inst.opCodes[j].Op
  })

  return s.writeFile(
    book,
    inst,
    "codes",
    "instructions",
    "Instruction List by name",
    "",
  )
}

func (s *M6502) writeOpsHexGrid(ctx context.Context) error {
  _ = s.autodoc.GetApi(ctx)

  book := generator.GetBook(ctx)
  inst := s.Instructions(book)

  return util.ReferenceFileBuilder("Opcode Matrix", "Instructions shown in an Opcode Matrix", "manual", 10).
      Then(NewHexGrid().
        Opcode(inst.opCodes...).
        FileBuilder()).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "hexgrid", "_index.html"), book.Modified())
}

func (s *M6502) writeFile(book *hugo.Book, inst *Instructions, prefix, name, title, desc string) error {
  err := util.ReferenceFileBuilder(
    title,
    desc,
    "manual",
    10,
  ).
    //Then(inst.writeOpCodes(prefix, inst.opCodes)).
    WrapAsFrontMatter().
      Then(func(slice strings2.StringSlice) (strings2.StringSlice, error) {
        slice = append(slice, "<div class='opIndex'>", "<table>")
        for _, op := range inst.opCodes {
          class := ""
          if op.Colour == "undocumented" {
            class = " class=\"" + op.Colour + "\""
          }
          slice = append(slice, fmt.Sprintf("<tr%s><td>%s</td><td>%s</td></tr>", class, op.String(), op.Code))
        }
        slice = append(slice, "</table>", "</div>")
        return slice, nil
      }).
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), name, "_index.html"), book.Modified())
  if err != nil {
    return err
  }

  // Use a copy as we reorder s.opCodes so when excel is processed later it gets just the last ordering
  codes := append([]*Opcode{}, inst.opCodes...)

  return util.WithTable().
    AsCSV(book.StaticPath(name+".csv"), book.Modified()).
    AsExcel(s.excel.Get(book.ID, book.Modified())).
      Do(&util.Table{
        Title: name,
        Columns: []string{
          "Decimal",
          "Hex",
          "Instruction",
          "Addressing",
          "Len(byte)",
          "Cycles",
        },
        RowCount: len(codes),
        GetRow: func(r int) interface{} {
          return codes[r]
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
      })
}
