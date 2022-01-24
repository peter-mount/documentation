package m6502

import (
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/generator/autodoc"
  util2 "github.com/peter-mount/go-kernel/util"
  "github.com/peter-mount/go-kernel/util/task"
)

type M6502 struct {
  generator    *generator.Generator `kernel:"inject"` // Generator
  excel        *generator.Excel     `kernel:"inject"` // Excel
  autodoc      *autodoc.Autodoc     `kernel:"inject"` // ResourceManager
  extracted    util2.Set            // Set of book ID's so that we run once per book
  instructions util2.Map            // Map of extracted data
}

func (s *M6502) Start() error {
  s.extracted = util2.NewHashSet()
  s.instructions = util2.NewSyncMap()

  s.generator.
      Register("6502OpsIndex",
        task.Of().
          Then(s.extractOpcodes).
          Then(delayOpTask(s.writeOpsIndex))).
      Register("6502OpsHexIndex",
        task.Of().
          Then(s.extractOpcodes).
          Then(delayOpTask(s.writeOpsHexIndex))).
      Register("6502OpsHexGrid",
        task.Of().
          Then(s.extractOpcodes).
          Then(delayOpTask(s.writeOpsHexGrid)))

  return nil
}
