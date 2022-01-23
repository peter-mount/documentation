package m6502

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/strings"
  "strconv"
  strings2 "strings"
)

const (
  InstructionsKey = "Instructions"
)

type Instructions struct {
  opCodes []*Opcode
  notes   *util.Notes
}

func GetInstructions(ctx context.Context) *Instructions {
  if i, exists := ctx.Value(InstructionsKey).(*Instructions); exists {
    return i
  }
  return nil
}

func (s *M6502) Instructions(b *hugo.Book) *Instructions {
  return s.instructions.ComputeIfAbsent(b.ID, func(key string) interface{} {
    return &Instructions{
      opCodes: nil,
      notes:   util.NewNotes(),
    }
  }).(*Instructions)
}

func (i *Instructions) decodeOpType(n *util.Notes, e1 interface{}) *OpcodeType {
  o := &OpcodeType{}

  if s, ok := e1.(string); ok {
    o.Value = s
  } else if i, ok := e1.(int); ok {
    o.Value = strconv.Itoa(i)
  } else {
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
  }

  return o
}

func (i *Instructions) extractOp(defaultOp string, n *util.Notes, e1 interface{}) {
  _ = util.IfMap(e1, func(e map[interface{}]interface{}) error {
    opcode := util.DecodeString(e["code"], "")
    order, _ := strconv.ParseInt(opcode, 16, 32)
    op := &Opcode{
      Order:         int(order),
      Code:          opcode,
      Op:            util.DecodeString(e["op"], defaultOp),
      Addressing:    util.DecodeString(e["addressing"], ""),
      Compatibility: util.NewSortedMap().Decode(e["compatibility"]),
      Colour:        util.DecodeString(e["colour"], ""),
    }

    if e1, exists := e["size"]; exists {
      op.Bytes = i.decodeOpType(n, e1)
    } else {
      op.Bytes = i.decodeOpType(n, e["bytes"])
    }

    if op.Bytes != nil && op.Bytes.Int() > 1 && !strings2.Contains(op.Code, "nn") {
      op.Code = op.Code + strings2.Repeat("nn", op.Bytes.Int()-1)
    }

    op.Cycles = i.decodeOpType(n, e["cycles"])

    i.opCodes = append(i.opCodes, op)

    return nil
  })
}

func (i *Instructions) writeOpCodes(prefix string, codes []*Opcode) util.FileBuilder {
  return func(a strings.StringSlice) (strings.StringSlice, error) {
    a = append(a, prefix+":")

    for _, op := range codes {
      a = append(a,
        "  - code: \""+op.Code+"\"",
        "    op: \""+op.Op+"\"",
        "    addressing: "+op.Addressing,
        "    compatibility:",
      )

      // Compatibility table is just the existence of the keys.
      // Sorted so we keep the same order each time
      _ = op.Compatibility.
        Keys().
        Sort().
          ForEach(func(k string) error {
            a = append(a, fmt.Sprintf("      %v: true", k))
            return nil
          })

      a = op.Bytes.append("    ", "bytes", a)
      a = op.Cycles.append("    ", "cycles", a)

      if op.Colour != "" {
        a = append(a, "    colour: \""+op.Colour+"\"")
      }
    }

    a = append(a, "notes:")
    for _, n := range i.notes.Notes {
      a = append(a, fmt.Sprintf("  - \"%s\" # %d", n.Value, n.Key))
    }
    return a, nil
  }
}
