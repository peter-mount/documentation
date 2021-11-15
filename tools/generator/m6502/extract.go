package m6502

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/walk"
  "log"
)

func (s *M6502) extractOpcodes(ctx context.Context) error {
  book := generator.GetBook(ctx)
  instructions := s.Instructions(book)

  // Only run once per Book ID
  if s.extracted.Contains(book.ID) {
    return nil
  }

  // Prevent us running again for this Book
  s.extracted.Add(book.ID)

  log.Println("Scanning 6502 opcodes")
  err := walk.NewPathWalker().
    IsFile().
    PathNotContain("/reference/").
    PathHasSuffix(".html").
      Then(hugo.FrontMatterActionOf().
        Then(s.extract).
        WithNotes(instructions.notes).
        Context(InstructionsKey, instructions).
        Walk(ctx)).
    Walk(book.ContentPath())
  if err != nil {
    return err
  }

  for _, op := range instructions.opCodes {
    op.Bytes.resolve(instructions.notes)
    op.Cycles.resolve(instructions.notes)
  }

  instructions.normalise()
  //instructions.validateOpcodes()

  return nil
}

func (s *M6502) extract(ctx context.Context, fm *hugo.FrontMatter) error {
  if codes, exists := fm.Other["codes"]; exists {
    var defaultOp string
    if a, exists := fm.Other["op"]; exists {
      defaultOp = a.(string)
    }

    notes := ctx.Value("notes").(*util.Notes)
    instructions := GetInstructions(ctx)

    _ = util.ForEachInterface(codes, func(e1 interface{}) error {
      instructions.extractOp(defaultOp, notes, e1)
      return nil
    })
  }
  return nil
}
