package m6502

import (
	"fmt"
	"github.com/peter-mount/area51.dev/area51.dev/tools/util"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Opcodes struct {
	Codes []Opcode `yaml:"codes"`
}

type Opcode struct {
	Order         int
	Code          string
	Op            string
	Addressing    string
	Compatibility map[interface{}]interface{}
	Bytes         *OpcodeType
	Cycles        *OpcodeType
	Notes         []int
}

type OpcodeType struct {
	Value  string       // value field
	NoteId []string     // Note id's
	Notes  []*util.Note // Resolved global note
}

func (o *OpcodeType) append(p, l string, a []string) []string {
	a = append(a,
		p+l+":",
		p+"  value: \""+o.Value+"\"",
	)

	if len(o.Notes) > 0 {
		a = append(a, p+"  notes:")
		for _, n := range o.Notes {
			a = append(a, p+"    - "+strconv.Itoa(n.Key))
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

func (s *M6502) normalise() {
	s.notes.Normalise()

	for _, o := range s.opCodes {
		o.Bytes.Normalise()
		o.Cycles.Normalise()
	}
}

func (s *M6502) extractOpcodes() error {
	s.opCodes = nil
	s.notes = util.NewNotes()

	log.Println("Scanning 6502 opcodes")
	err := filepath.Walk(*s.baseDir, s.extractOpcode)
	if err != nil {
		return err
	}

	for _, op := range s.opCodes {
		op.Bytes.resolve(s.notes)
		op.Cycles.resolve(s.notes)
	}

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

	if e, ok := e1.(map[interface{}]interface{}); ok {
		if v, ok := e["value"]; ok {
			o.Value = util.DecodeString(v, "")
		}

		if v, ok := e["notes"]; ok {
			if a, ok := v.([]interface{}); ok {
				for _, ae := range a {
					if i, ok := ae.(int); ok {
						note := n.GetId(i)
						if note != nil {
							o.NoteId = append(o.NoteId, note.Value)
						}
					}
				}
			}
		}
	}

	return o
}

func (s *M6502) extractOp(defaultOp string, n *util.Notes, e1 interface{}) {

	e, ok := e1.(map[interface{}]interface{})
	if !ok {
		log.Println(ok, e1)
		return
	} else {
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
	}
}

func (s *M6502) writeOpsIndex() error {
	sort.SliceStable(s.opCodes, func(i, j int) bool {
		return s.opCodes[i].Order < s.opCodes[j].Order
	})

	return s.writeFile(
		"opcodes",
		"Instruction List by opcode",
		"6502 instructions by hex opcode",
	)
}

func (s *M6502) writeOpsHexIndex() error {
	sort.SliceStable(s.opCodes, func(i, j int) bool {
		if s.opCodes[i].Op == s.opCodes[j].Op {
			return s.opCodes[i].Addressing < s.opCodes[j].Addressing
		}
		return s.opCodes[i].Op < s.opCodes[j].Op
	})

	return s.writeFile(
		"instructions",
		"Instruction List by name",
		"6502 instructions by name",
	)
}

func (s *M6502) writeFile(name, title, desc string) error {
	return util.WriteReferenceFile(
		*s.baseDir,
		name,
		title,
		desc,
		10,
		func(a []string) ([]string, error) {
			a = append(a, "codes:")

			for _, op := range s.opCodes {
				a = append(a,
					"  - code: \""+op.Code+"\"",
					"    op: \""+op.Op+"\"",
					"    addressing: "+op.Addressing,
					"    compatibility:",
				)

				for k, _ := range op.Compatibility {
					a = append(a, fmt.Sprintf("      %v: true", k))
				}

				a = op.Bytes.append("    ", "bytes", a)
				a = op.Cycles.append("    ", "cycles", a)
			}

			a = append(a, "notes:")
			for _, n := range s.notes.Notes {
				a = append(a, fmt.Sprintf("  - \"%s\"", n.Value))
			}
			return a, nil
		})
}

func (s *M6502) appendMap(p, l string, a []string, m map[interface{}]interface{}) []string {
	a = append(a, p+l+":")

	for kv, v := range m {
		k := kv.(string)
		if v == nil {

		} else if s, ok := v.(string); ok {
			a = append(a, fmt.Sprintf("%s  %s: \"%s\"", p, k, s))
		} else if s, ok := v.(int); ok {
			a = append(a, fmt.Sprintf("%s  %s: %d", p, k, s))
		} else if s, ok := v.([]interface{}); ok {
			a = append(a, fmt.Sprintf("%s  %s:", p, k))
			for _, e := range s {
				a = append(a, fmt.Sprintf("%s    - %v", p, e))
			}
		} else {
			log.Printf("k=%v %T v=%v", kv, v, v)
		}
	}

	return a
}
