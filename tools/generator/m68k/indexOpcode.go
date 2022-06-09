package m68k

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/generator"
  strings2 "github.com/peter-mount/go-kernel/util/strings"
)

func (s *M68k) writeOpcodeIndex(ctx context.Context) error {
  _ = s.autodoc.GetApi(ctx)

  book := generator.GetBook(ctx)
  inst := s.Instructions(book).
      Sort(func(a *Opcode, b *Opcode) bool {
        return a.Order < b.Order
      })

  gen := &IndexGenerator{
    Prefix:    "codes",
    Name:      "opcodes",
    Title:     "Instruction List by Opcode",
    Desc:      "",
    Class:     "opIndex2",
    Paginator: nil,
    Header:    nil,
    Body: func(slice strings2.StringSlice, _ int, entry interface{}) strings2.StringSlice {
      op := entry.(*Opcode)
      class := ""
      if op.Colour == "undocumented" {
        class = " class=\"" + op.Colour + "\""
      }
      return append(slice, fmt.Sprintf("<tr%s><td>%s</td><td>%s</td></tr>", class, op.String(), op.Code))
    },
  }
  return gen.WriteFile(book, inst.Iterator())
}
