package assembly

import (
	"context"
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/documentation/tools/gensite/util"
	util2 "github.com/peter-mount/go-kernel/v2/util"
	"sort"
	"strconv"
	strings2 "strings"
)

const (
	InstructionsKey = "Instructions"
)

// Instructions holds the instructions being processed
type Instructions struct {
	opCodes         []*Opcode                       // Slice of Opcodes
	notes           *util.Notes                     // Notes, if any
	ExtractFormat   func(format string, op *Opcode) // Extract format field (optional)
	OpcodeFormatter func(op *Opcode) string         // Format Opcode
}

func NewInstructions() *Instructions {
	return &Instructions{
		opCodes:         nil,
		notes:           util.NewNotes(),
		OpcodeFormatter: defaultOpcodeFormatter,
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
func (i *Instructions) Extract(defaultOp string, n *util.Notes, e1 interface{}) {
	_ = util2.IfMap(e1, func(e map[interface{}]interface{}) error {

		op := &Opcode{
			Code:          util2.DecodeString(e["code"], ""),
			Op:            util2.DecodeString(e["op"], defaultOp),
			Addressing:    util2.DecodeString(e["addressing"], ""),
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

func (i *Instructions) ExtractFrontMatter(ctx context.Context, fm *hugo.FrontMatter) error {
	if codes, exists := fm.Other["codes"]; exists {
		var defaultOp string
		if a, exists := fm.Other["op"]; exists {
			defaultOp = a.(string)
		}

		notes := ctx.Value("notes").(*util.Notes)

		_ = util2.ForEachInterface(codes, func(e1 interface{}) error {
			i.Extract(defaultOp, notes, e1)
			return nil
		})
	}
	return nil
}

func (i *Instructions) decodeOpType(n *util.Notes, e1 interface{}) *OpcodeType {
	o := &OpcodeType{}

	if s, ok := e1.(string); ok {
		o.Value = s
	} else if i, ok := e1.(int); ok {
		o.Value = strconv.Itoa(i)
	} else {
		_ = util2.IfMap(e1, func(e map[interface{}]interface{}) error {
			_ = util2.IfMapEntry(e, "value", func(v interface{}) error {
				o.Value = util2.DecodeString(v, "")
				return nil
			})

			return util2.IfMapEntry(e, "notes", func(v interface{}) error {
				return util2.ForEachInterface(v, func(ae interface{}) error {
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

func (i *Instructions) Normalise() {
	for _, op := range i.opCodes {
		op.Bytes.resolve(i.notes)
		op.Cycles.resolve(i.notes)
	}

	i.notes.Normalise()
}

func defaultOpcodeFormatter(op *Opcode) string {
	return op.Op
}
