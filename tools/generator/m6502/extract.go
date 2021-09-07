package m6502

import (
  "context"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/walk"
  "log"
  "strconv"
)

func (s *M6502) extractOpcodes(book *hugo.Book) error {
  s.opCodes = nil
  s.notes = util.NewNotes()

  log.Println("Scanning 6502 opcodes")
  err := walk.NewPathWalker().
    IsFile().
    PathNotContain("/reference/").
    PathHasSuffix(".html").
      Then(hugo.FrontMatterActionOf().
        Then(s.extract).
        WithNotes(s.notes).
        Walk()).
    Walk(book.ContentPath())
  if err != nil {
    return err
  }

  for _, op := range s.opCodes {
    op.Bytes.resolve(s.notes)
    op.Cycles.resolve(s.notes)
  }

  s.normalise()

  s.validateOpcodes()

  s.extracted = true
  return nil
}

func (s *M6502) extract(ctx context.Context, fm *hugo.FrontMatter) error {
  if codes, exists := fm.Other["codes"]; exists {
    var defaultOp string
    if a, exists := fm.Other["op"]; exists {
      defaultOp = a.(string)
    }

    notes := ctx.Value("notes").(*util.Notes)

    _ = util.ForEachInterface(codes, func(e1 interface{}) error {
      s.extractOp(defaultOp, notes, e1)
      return nil
    })
  }
  return nil
}

func (s *M6502) decodeOpType(n *util.Notes, e1 interface{}) *OpcodeType {
  o := &OpcodeType{}

  _ = util.IfMap(e1, func(e map[interface{}]interface{}) error {
    _ = util.IfMapEntry(e, "value", func(v interface{}) error {
      o.Value = util.DecodeString(v, "")
      return nil
    })

    return util.IfMapEntry(e, "notes", func(v interface{}) error {
      return util.ForEachInterface(v, func(ae interface{}) error {
        if i, ok := ae.(int); ok {
          note := n.GetId(i)
          if note != nil {
            o.NoteId = append(o.NoteId, note.Value)
          }
        }
        return nil
      })
    })
  })

  return o
}

func (s *M6502) extractOp(defaultOp string, n *util.Notes, e1 interface{}) {
  _ = util.IfMap(e1, func(e map[interface{}]interface{}) error {
    opcode := util.DecodeString(e["code"], "")
    order, _ := strconv.ParseInt(opcode, 16, 32)
    op := &Opcode{
      Order:         int(order),
      Code:          opcode,
      Op:            util.DecodeString(e["op"], defaultOp),
      Addressing:    util.DecodeString(e["addressing"], ""),
      Compatibility: e["compatibility"].(map[interface{}]interface{}),
    }

    op.Bytes = s.decodeOpType(n, e["bytes"])
    op.Cycles = s.decodeOpType(n, e["cycles"])

    s.opCodes = append(s.opCodes, op)

    return nil
  })
}
