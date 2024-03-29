package m68k

import (
	"github.com/peter-mount/documentation/tools/gensite/generator"
	"github.com/peter-mount/documentation/tools/gensite/generator/assembly"
	"github.com/peter-mount/documentation/tools/gensite/generator/autodoc"
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

type M68k struct {
	generator    *generator.Generator `kernel:"inject"` // Generator
	excel        *generator.Excel     `kernel:"inject"` // Excel
	autodoc      *autodoc.Autodoc     `kernel:"inject"` // ResourceManager
	extracted    util.Set             // Set of book ID's so that we run once per book
	instructions util.Map             // Map of extracted data
}

func (s *M68k) Start() error {
	s.extracted = util.NewHashSet()
	s.instructions = util.NewSyncMap()

	s.generator.
		Register("68kOperationIndex",
			task.Of().
				Then(s.extractOpcodes).
				Then(assembly.DelayOpTask(s.writeOperationIndex)).
				Then(assembly.DelayOpTask(s.writeOpcodeIndex)))

	return nil
}

func (s *M68k) Instructions(b *hugo.Book) *assembly.Instructions {
	return s.instructions.ComputeIfAbsent(b.ID, assembly.ComputeNewInstructions).(*assembly.Instructions)
}
