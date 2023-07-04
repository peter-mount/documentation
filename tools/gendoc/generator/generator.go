package generator

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/gendoc"
	"github.com/peter-mount/documentation/tools/gendoc/hugo"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

// Generator is a kernel Service which handles the generation of content based on page metadata.
type Generator struct {
	bookShelf  *hugo.BookShelf      `kernel:"inject"`
	worker     task.Queue           `kernel:"worker"` // Worker queue
	generators map[string]task.Task // Map of available generators
}

func (g *Generator) Start() error {
	g.generators = make(map[string]task.Task)

	g.worker.AddPriorityTask(tools.PriorityImmediate, g.generate)
	return nil
}

// Register a named Handler
func (g *Generator) Register(n string, h task.Task) *Generator {
	if _, exists := g.generators[n]; exists {
		panic(fmt.Errorf("GeneratorHandler %s already registered", n))
	}

	g.generators[n] = h
	return g
}

func (g *Generator) generate(_ context.Context) error {
	if err := g.bookShelf.Books().ForEach(g.invokeBook); err != nil {
		return err
	}
	return nil
}

const (
	BookKey = "hugo.Book"
)

func GetBook(ctx context.Context) *hugo.Book {
	if b, ok := ctx.Value(BookKey).(*hugo.Book); ok {
		return b
	}
	return nil
}

func (g *Generator) invokeBook(book *hugo.Book) error {
	return book.Generate.ForEach(func(n string) error {
		h, exists := g.generators[n]
		if exists {
			g.worker.AddTask(h.WithValue(BookKey, book))
		} else {
			// Log a warning but ignore - could be an invalid config or the generator is not deployed.
			// Originally this was a fatal error, but now we just ignore to allow custom tools to be run
			log.Printf("book %s GeneratorHandler %s is not registered", book.ID, n)
		}

		return nil
	})
}
