package m68k

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/generator/assembly"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  util2 "github.com/peter-mount/go-kernel/util"
  "github.com/peter-mount/go-kernel/util/walk"
  "log"
  "strconv"
)

func (s *M68k) extractOpcodes(ctx context.Context) error {
  book := generator.GetBook(ctx)

  instructions := s.Instructions(book)
  instructions.ExtractFormat = s.extractFormat

  // Only run once per Book ID
  if s.extracted.Contains(book.ID) {
    return nil
  }

  // Prevent us running again for this Book
  s.extracted.Add(book.ID)

  log.Println("Scanning 68K opcodes")
  err := walk.NewPathWalker().
    IsFile().
    PathNotContain("/reference/").
    PathHasSuffix(".html").
      Then(hugo.FrontMatterActionOf().
        Then(s.extract).
        WithNotes(instructions.Notes()).
        Context(assembly.InstructionsKey, instructions).
        Walk(ctx)).
    Walk(book.ContentPath())
  if err != nil {
    return err
  }

  instructions.Normalise()

  return nil
}

func (s *M68k) extract(ctx context.Context, fm *hugo.FrontMatter) error {
  if codes, exists := fm.Other["codes"]; exists {
    var defaultOp string
    if a, exists := fm.Other["op"]; exists {
      defaultOp = a.(string)
    }

    notes := ctx.Value("notes").(*util.Notes)
    instructions := assembly.GetInstructions(ctx)

    _ = util2.ForEachInterface(codes, func(e1 interface{}) error {
      instructions.Extract(defaultOp, notes, e1)
      return nil
    })
  }
  return nil
}

func (s *M68k) extractFormat(format string, op *assembly.Opcode) {

  var opcodeR []rune
  var orderR []rune
  for _, c := range format {
    cc := 'X'
    cr := '0'
    cn := 1
    switch c {
    case '0':
      cc = '0'
      cr = '0'
    case '1':
      cc = '1'
      cr = '1'
    case 'd':
      cn = 3
    case 'E':
      cn = 6
    case 'M':
      cn = 1
    case 'm':
      cn = 3
    case 'R':
      cn = 3
    case 'r':
      cn = 3
    case 'S':
      cn = 2
    }
    for i := 0; i < cn; i++ {
      opcodeR = append(opcodeR, cc)
      orderR = append(orderR, cr)
    }
  }

  op.Code = string(opcodeR)
  order, _ := strconv.ParseInt(string(orderR), 2, 32)
  op.Order = int(order)
}
