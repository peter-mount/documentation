package autodoc

import "github.com/peter-mount/go-kernel/v2/util"

type FunctionParams struct {
	A string // A register
	X string // X register
	Y string // Y register
	C string // C Carry Flag
}

func (p *FunctionParams) decode(e interface{}) error {
	return util.IfMap(e, func(m map[interface{}]interface{}) error {
		p.A = util.IfMapEntryString(m, "a")
		p.X = util.IfMapEntryString(m, "x")
		p.Y = util.IfMapEntryString(m, "y")
		p.C = util.IfMapEntryString(m, "c")
		return nil
	})
}
