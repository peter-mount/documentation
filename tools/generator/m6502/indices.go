package m6502

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/generator"
	"github.com/peter-mount/documentation/tools/generator/assembly"
	strings2 "github.com/peter-mount/go-kernel/v2/util/strings"
)

func (s *M6502) writeOpsIndex(ctx context.Context) error {
	_ = s.autodoc.GetApi(ctx)

	book := generator.GetBook(ctx)
	inst := s.Instructions(book).
		Sort(func(a *assembly.Opcode, b *assembly.Opcode) bool {
			return assembly.DecodeOpcode(a.Code) < assembly.DecodeOpcode(b.Code)
		})

	gen := &assembly.IndexGenerator{
		Prefix:    "codes",
		Name:      "opcodes",
		Title:     "Instruction List by opcode",
		Desc:      "",
		Class:     "opIndex",
		Paginator: nil,
		Header:    nil,
		Body:      m6502IndexBody(inst),
	}
	return gen.WriteFile(book, inst.Iterator())
}

func (s *M6502) writeOpsHexIndex(ctx context.Context) error {
	_ = s.autodoc.GetApi(ctx)

	book := generator.GetBook(ctx)
	inst := s.Instructions(book).
		Sort(func(a *assembly.Opcode, b *assembly.Opcode) bool {
			if a.Op == b.Op {
				return a.Addressing < b.Addressing
			}
			return a.Op < b.Op
		})

	gen := &assembly.IndexGenerator{
		Prefix:    "codes",
		Name:      "instructions",
		Title:     "Instruction List by name",
		Desc:      "",
		Class:     "opIndex",
		Paginator: nil,
		Header:    nil,
		Body:      m6502IndexBody(inst),
	}
	return gen.WriteFile(book, inst.Iterator())
}

// m6502IndexBody shared
func m6502IndexBody(inst *assembly.Instructions) func(strings2.StringSlice, int, interface{}) strings2.StringSlice {
	return func(slice strings2.StringSlice, _ int, entry interface{}) strings2.StringSlice {
		op := entry.(*assembly.Opcode)
		class := ""
		if op.Colour == "undocumented" {
			class = " class=\"" + op.Colour + "\""
		}
		return append(slice, fmt.Sprintf("<tr%s><td>%s</td><td>%s</td></tr>", class, inst.OpcodeFormatter(op), op.Code))
	}
}
