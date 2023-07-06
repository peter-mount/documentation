package m6502

import (
	"context"
	"github.com/peter-mount/documentation/tools/gensite/generator"
	"github.com/peter-mount/documentation/tools/gensite/generator/assembly"
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util/walk"
)

func (s *M6502) extractOpcodes(ctx context.Context) error {
	book := generator.GetBook(ctx)
	instructions := s.Instructions(book)

	// Only run once per Book ID
	if s.extracted.Contains(book.ID) {
		return nil
	}

	// Prevent us running again for this Book
	s.extracted.Add(book.ID)

	log.Println("Scanning 6502 opcodes")
	err := walk.NewPathWalker().
		IsFile().
		PathNotContain("/reference/").
		PathHasSuffix(".html").
		Then(hugo.FrontMatterActionOf().
			Then(instructions.ExtractFrontMatter).
			WithNotes(instructions.Notes()).
			Context(assembly.InstructionsKey, instructions).
			Walk(ctx)).
		Walk(book.ContentPath())
	if err != nil {
		return err
	}

	instructions.Normalise()

	return nil
}
