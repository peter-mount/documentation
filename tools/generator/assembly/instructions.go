package assembly

import (
  "context"
  "github.com/peter-mount/documentation/tools/util"
  util2 "github.com/peter-mount/go-kernel/util"
  "sort"
  "strconv"
)

const (
  InstructionsKey = "Instructions"
)

// Instructions holds the instructions being processed
type Instructions struct {
  opCodes       []*Opcode                       // Slice of Opcodes
  notes         *util.Notes                     // Notes, if any
  ExtractFormat func(format string, op *Opcode) // Extract format field (optional)
}

func NewInstructions() *Instructions {
  return &Instructions{
    opCodes: nil,
    notes:   util.NewNotes(),
  }
}

func ComputeNewInstructions(_ string) interface{} {
  return NewInstructions()
}

func (i *Instructions) Notes() *util.Notes {
  return i.notes
}

func (i *Instructions) Sort(comparator func(a *Opcode, b *Opcode) bool) *Instructions {
  sort.SliceStable(i.opCodes, func(a, b int) bool {
    return comparator(i.opCodes[a], i.opCodes[b])
  })
  return i
}

func (i *Instructions) Iterator() util2.Iterator {
  return &iterator{
    src: i.opCodes,
    i:   0,
    l:   len(i.opCodes),
  }
}

type iterator struct {
  src  []*Opcode
  i, l int
}

func (i *iterator) HasNext() bool {
  return i.i < i.l
}

func (i *iterator) Next() interface{} {
  v := i.src[i.i]
  i.i++
  return v
}

func (i *iterator) ForEach(f func(interface{})) {
  for _, v := range i.src {
    f(v)
  }
}

func (i *iterator) ForEachAsync(f func(interface{})) {
  i.ForEach(f)
}

func (i *iterator) ForEachFailFast(f func(interface{}) error) error {
  for _, v := range i.src {
    err := f(v)
    if err != nil {
      return err
    }
  }
  return nil
}

func (i *iterator) Iterator() util2.Iterator {
  // Just return ourselves
  return i
}

func (i *iterator) ReverseIterator() util2.Iterator {
  // Just return ourselves
  return i
}

func GetInstructions(ctx context.Context) *Instructions {
  if i, exists := ctx.Value(InstructionsKey).(*Instructions); exists {
    return i
  }
  return nil
}

// Extract op code.
// TODO this is 68k specific
func (i *Instructions) Extract(defaultOp string, n *util.Notes, e1 interface{}) {
  _ = util2.IfMap(e1, func(e map[interface{}]interface{}) error {

    op := &Opcode{
      Code:          util2.DecodeString(e["code"], ""),
      Op:            util2.DecodeString(e["op"], defaultOp),
      Format:        util2.DecodeString(e["format"], ""),
      Compatibility: util2.NewSortedMap().Decode(e["compatibility"]),
      Colour:        util2.DecodeString(e["colour"], ""),
    }

    order, _ := strconv.ParseInt(op.Code, 16, 32)
    op.Order = int(order)

    if i.ExtractFormat != nil {
      if op.Format != "" {
        i.ExtractFormat(op.Format, op)
      }
    }

    i.opCodes = append(i.opCodes, op)

    return nil
  })
}

func (i *Instructions) Normalise() {
  for _, op := range i.opCodes {
    op.Bytes.resolve(i.notes)
    op.Cycles.resolve(i.notes)
  }

  i.notes.Normalise()
}
