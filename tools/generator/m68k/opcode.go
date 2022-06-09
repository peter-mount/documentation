package m68k

import (
  "github.com/peter-mount/documentation/tools/util"
  util2 "github.com/peter-mount/go-kernel/util"
  "sort"
  "strconv"
)

type Opcodes struct {
  Codes []Opcode `yaml:"codes"`
}

type Opcode struct {
  Order         int              // numerical order of the Opcode, usually same as Code
  Code          string           // hex code of Opcode as a string
  Op            string           // Operation name
  Format        string           // Binary format of Opcode (optional)
  Compatibility *util2.SortedMap // Opcode compatibility map
  Notes         []int            // Notes about opcode
  Colour        string           // Colour used in rendering (optional)
}

// String returns the Op field
func (op *Opcode) String() string {
  return op.Op
}

// OpcodeType handles details about an Opcode
type OpcodeType struct {
  Value  string       // value field
  NoteId []string     // Note id's
  Notes  []*util.Note // Resolved global note
}

// String returns Value field
func (o *OpcodeType) String() string {
  if o != nil {
    return o.Value
  }
  return ""
}

// Int returns Value as an int, 0 if invalid
func (o *OpcodeType) Int() int {
  if o != nil {
    i, err := strconv.Atoi(o.Value)
    if err == nil {
      return i
    }
  }
  return 0
}

func (o *OpcodeType) append(p, l string, a []string) []string {
  a = append(a,
    p+l+":",
    p+"  value: \""+o.Value+"\"",
  )

  if len(o.Notes) > 0 {
    a = append(a, p+"  notes:")
    for _, n := range o.Notes {
      a = append(a, p+"    - "+strconv.Itoa(n.Key)+" # "+n.Value)
    }
  }
  return a
}

func (o *OpcodeType) resolve(n *util.Notes) {
  if o == nil {
    return
  }

  for _, s := range o.NoteId {
    note := n.Get(s)
    if note != nil {
      o.Notes = append(o.Notes, note)
    }
  }
}

func (o *OpcodeType) Normalise() {
  sort.SliceStable(o.Notes, func(i, j int) bool {
    return o.Notes[i].Key < o.Notes[j].Key
  })
}

func (i *Instructions) normalise() {
  i.notes.Normalise()
}
