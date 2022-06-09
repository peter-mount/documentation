package m6502

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/util"
)

func (s *M6502) writeOpsHexGrid(ctx context.Context) error {
  _ = s.autodoc.GetApi(ctx)

  book := generator.GetBook(ctx)
  inst := s.Instructions(book)

  return util.ReferenceFileBuilder("Opcode Matrix", "Instructions shown in an Opcode Matrix", "manual", 10, book.Modified()).
      Then(NewHexGrid().
        OpcodeFrom(inst.Iterator()).
        FileBuilder()).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "hexgrid", "_index.html"), book.Modified())
}
