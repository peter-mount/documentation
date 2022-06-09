package m68k

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/generator/assembly"
  strings2 "github.com/peter-mount/go-kernel/util/strings"
)

func (s *M68k) writeOperationIndex(ctx context.Context) error {
  _ = s.autodoc.GetApi(ctx)

  book := generator.GetBook(ctx)
  inst := s.Instructions(book).
      Sort(func(a *assembly.Opcode, b *assembly.Opcode) bool {
        return decodeOpcode(a.Code) < decodeOpcode(b.Code)
      })

  gen := &assembly.IndexGenerator{
    Prefix:    "codes",
    Name:      "operation",
    Title:     "Instruction List by Operation",
    Desc:      "",
    Class:     "opIndex1",
    Paginator: nil,
    Header: func(slice strings2.StringSlice, rowCount int) strings2.StringSlice {
      return append(slice, "<thead><tr>",
        "<th>Operation</th>",
        "<th>Opcode</th>",
        "</tr></thead>")
    },
    Body: func(slice strings2.StringSlice, _ int, entry interface{}) strings2.StringSlice {
      op := entry.(*assembly.Opcode)
      class := ""
      if op.Colour == "undocumented" {
        class = " class=\"" + op.Colour + "\""
      }
      return append(slice, fmt.Sprintf("<tr%s><td>%s</td><td>%s</td></tr>", class, op.String(), op.Code))
    },
  }
  return gen.WriteFile(book, inst.Iterator())
}
