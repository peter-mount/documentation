package m6502

import (
  "github.com/peter-mount/documentation/tools/util"
  "sort"
  "strconv"
)

type Opcodes struct {
  Codes []Opcode `yaml:"codes"`
}

type Opcode struct {
  Order         int
  Code          string
  Op            string
  Addressing    string
  Compatibility *util.SortedMap
  Bytes         *OpcodeType
  Cycles        *OpcodeType
  Notes         []int
  Colour        string
}

func (op *Opcode) String() string {
  switch op.Addressing {
  case "abs":
    return op.Op + " <em>addr</em>"
  case "absi":
    return op.Op + " (<em>addr</em>)"
  case "absii":
    return op.Op + " (<em>addr</em>,X)"
  case "absix":
    return op.Op + " <em>addr</em>,X"
  case "absiy":
    return op.Op + " <em>addr</em>,Y"
  case "abslix":
    return op.Op + " <em>long</em>,X"
  case "absl":
    return op.Op + " <em>long</em>"
  case "absil":
    return op.Op + " [<em>addr</em>]"
  case "acc":
    return op.Op + " A"
  case "bm":
    return op.Op + " <em>srcbk</em>, <em>dstbk</em>"
  case "dp":
    return op.Op + " <em>dp</em>"
  case "dpi":
    return op.Op + " (<em>dp</em>)"
  case "dpil":
    return op.Op + " [<em>dp</em>]"
  case "dpix":
    return op.Op + " <em>dp</em>,X"
  case "dpiix":
    return op.Op + " (<em>dp</em>,X)"
  case "dpiy":
    return op.Op + " <em>dp</em>,Y"
  case "dpiiy":
    return op.Op + " (<em>dp</em>),Y"
  case "dpiliy":
    return op.Op + " [<em>dp</em>],Y"
  case "imm":
    return op.Op + " #<em>const</em>"
  case "pcr":
    return op.Op + " <em>nearlabel</em>"
  case "pcrl":
    return op.Op + " <em>label</em>"
  case "sa":
    return op.Op + " <em>addr</em>"
  case "sdpi":
    return op.Op + " (<em>dp</em>)"
  case "sic":
    return op.Op + " <em>const</em>"
  case "spcrl":
    return op.Op + " <em>label</em>"
  case "sr":
    return op.Op + " <em>sr</em>,S"
  case "sriiy":
    return op.Op + " (<em>sr</em>,S),Y"
  default:
    return op.Op
  }
}

type OpcodeType struct {
  Value  string       // value field
  NoteId []string     // Note id's
  Notes  []*util.Note // Resolved global note
}

func (o *OpcodeType) String() string {
  if o != nil {
    return o.Value
  }
  return ""
}

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

  for _, o := range i.opCodes {
    o.Bytes.Normalise()
    o.Cycles.Normalise()
  }
}
