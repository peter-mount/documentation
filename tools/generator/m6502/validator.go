package m6502

import "log"

// A simple validator - checks the data for duplicate or missing opcodes
func (i *Instructions) validateOpcodes() {
	log.Println("Validation Opcode allocation")

	m := make(map[int]*Opcode)

	for _, o := range i.opCodes {
		if e, exists := m[o.Order]; exists {
			log.Printf("Duplicate entry for %03d %s %s - duplicate is %s",
				e.Order, e.Code,
				e.Op,
				o.Op,
			)
		} else {
			m[o.Order] = o
		}
	}

	for i := 0; i < 256; i++ {
		if _, exists := m[i]; !exists {
			log.Printf("Missing entry for %03d %02x", i, i)
		}
	}
}
