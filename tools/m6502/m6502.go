package m6502

import (
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
)

type M6502 struct {
  generator *hugo.Generator // Generator
  extracted bool            // true once extract has run
  opCodes   []*Opcode
  notes     *util.Notes
}

func (s *M6502) Name() string {
  return "m6502"
}

func (s *M6502) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&hugo.Generator{})
  if err != nil {
    return err
  }
  s.generator = service.(*hugo.Generator)

  return nil
}

func (s *M6502) Start() error {
  s.generator.
    RegisterCompound("6502OpsIndex", s.extractOpcodes, s.writeOpsIndex).
    RegisterCompound("6502OpsHexIndex", s.extractOpcodes, s.writeOpsHexIndex)

  return nil
}
