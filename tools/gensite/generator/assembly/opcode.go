package assembly

import (
	"github.com/peter-mount/go-kernel/v2/util"
)

type Opcode struct {
	Order         int                   // numerical order of the Opcode, usually same as Code
	Code          string                // hex code of Opcode as a string
	Op            string                // Operation name
	Addressing    string                // Addressing mode (optional)
	Format        string                // Binary format of Opcode (optional)
	Compatibility *util.SortedMap[bool] // Opcode compatibility map
	Bytes         *OpcodeType           // Bytes opcode uses (optional)
	Cycles        *OpcodeType           // Cycles opcode uses (optional)
	Notes         []int                 // Notes about opcode
	Colour        string                // Colour used in rendering (optional)
}
