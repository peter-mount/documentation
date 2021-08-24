package m6502

import (
	"flag"
	"github.com/peter-mount/documentation/tools/util"
	"github.com/peter-mount/go-kernel"
)

type M6502 struct {
	baseDir *string
	opCodes []*Opcode
	notes   *util.Notes
}

func (s *M6502) Name() string {
	return "m6502"
}

func (s *M6502) Init(_ *kernel.Kernel) error {
	s.baseDir = flag.String("6502", "", "Src dir for m6502 content")

	return nil
}

func (s *M6502) Run() error {
	// Do nothing
	if *s.baseDir == "" {
		return nil
	}

	err := s.extractOpcodes()
	if err != nil {
		return err
	}

	s.normalise()

	s.validateOpcodes()

	err = s.writeOpsIndex()
	if err != nil {
		return err
	}

	err = s.writeOpsHexIndex()
	if err != nil {
		return err
	}

	return nil
}
