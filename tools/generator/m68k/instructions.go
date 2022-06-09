package m68k

import (
  "context"
  "github.com/peter-mount/documentation/tools/hugo"
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
  opCodes []*Opcode
  notes   *util.Notes
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

func (s *M68k) Instructions(b *hugo.Book) *Instructions {
  return s.instructions.ComputeIfAbsent(b.ID, func(key string) interface{} {
    return &Instructions{
      opCodes: nil,
      notes:   util.NewNotes(),
    }
  }).(*Instructions)
}

func (i *Instructions) extractOp(defaultOp string, n *util.Notes, e1 interface{}) {
  _ = util2.IfMap(e1, func(e map[interface{}]interface{}) error {

    format := util2.DecodeString(e["format"], "")

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

    opcode := string(opcodeR)
    order, _ := strconv.ParseInt(string(orderR), 2, 32)
    op := &Opcode{
      Order:         int(order),
      Code:          opcode,
      Op:            util2.DecodeString(e["op"], defaultOp),
      Format:        format,
      Compatibility: util2.NewSortedMap().Decode(e["compatibility"]),
      Colour:        util2.DecodeString(e["colour"], ""),
    }

    i.opCodes = append(i.opCodes, op)

    return nil
  })
}
