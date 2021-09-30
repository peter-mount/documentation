package m6502

import (
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
)

type M6502 struct {
  generator *generator.Generator // Generator
  extracted bool                 // true once extract has run
  opCodes   []*Opcode
  notes     *util.Notes
}

func (s *M6502) Name() string {
  return "m6502"
}

func (s *M6502) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&generator.Generator{})
  if err != nil {
    return err
  }
  s.generator = service.(*generator.Generator)

  return nil
}

func (s *M6502) Start() error {
  s.generator.
      Register("6502OpsIndex",
        generator.HandlerOf().
          RunOnce(&s.extracted, s.extractOpcodes).
          Then(s.writeOpsIndex)).
      Register("6502OpsHexIndex",
        generator.HandlerOf().
          RunOnce(&s.extracted, s.extractOpcodes).
          Then(s.writeOpsHexIndex))

  return nil
}
