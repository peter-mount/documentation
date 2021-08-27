package m6502

import (
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "log"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

func (s *M6502) extractOpcodes(book *hugo.Book) error {
  if s.extracted {
    return nil
  }

  s.opCodes = nil
  s.notes = util.NewNotes()

  log.Println("Scanning 6502 opcodes")
  err := filepath.Walk(book.ContentPath(), s.extractOpcode)
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

func (s *M6502) extractOpcode(path string, info os.FileInfo, err error) error {
  if err != nil {
    return err
  }

  if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") && !strings.Contains(path, "reference") {
    //log.Println("Reading", path)
    m, err := util.LoadFrontMatter("", path)
    if err != nil {
      return err
    }

    notes := util.NewNotes()
    notes.DecodePageNotes(m["notes"])

    if codes, ok := m["codes"]; ok {
      var defaultOp string
      a, ok := m["op"]
      if ok {
        defaultOp = a.(string)
      }

      for _, e1 := range codes.([]interface{}) {
        s.extractOp(defaultOp, notes, e1)
      }
    }

    // Import these notes into the global pool
    s.notes.Merge(notes)
  }

  return nil
}

func (s *M6502) decodeOpType(n *util.Notes, e1 interface{}) *OpcodeType {
  o := &OpcodeType{}

  util.IfMap(e1, func(e map[interface{}]interface{}) {
    util.IfMapEntry(e, "entry", func(v interface{}) {
      o.Value = util.DecodeString(v, "")
    })

    util.IfMapEntry(e, "notes", func(v interface{}) {
      util.ForEachInterface(v, func(ae interface{}) {
        if i, ok := ae.(int); ok {
          note := n.GetId(i)
          if note != nil {
            o.NoteId = append(o.NoteId, note.Value)
          }
        }
      })
    })
  })

  return o
}

func (s *M6502) extractOp(defaultOp string, n *util.Notes, e1 interface{}) {
  util.IfMap(e1, func(e map[interface{}]interface{}) {
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
  })
}
