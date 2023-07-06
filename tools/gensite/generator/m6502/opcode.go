package m6502

import "github.com/peter-mount/documentation/tools/gensite/generator/assembly"

func M6502OpcodeFormatter(op *assembly.Opcode) string {
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
