package generator

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "log"
)

// GeneratorHandler performs an action against a Book
type GeneratorHandler func(*hugo.Book) error

// Then returns a GeneratorHandler that calls this one then if no error the next one forming a chain
func (h GeneratorHandler) Then(next GeneratorHandler) GeneratorHandler {
  return func(b *hugo.Book) error {
    err := h(b)
    if err != nil {
      return err
    }
    return next(b)
  }
}

// GeneratorHandlerOf returns a GeneratorHandler that will invoke each supplied GeneratorHandler in a single instance.
// This is a convenience form of calling Then() on each one.
func GeneratorHandlerOf(handlers ...GeneratorHandler) GeneratorHandler {
  switch len(handlers) {
  case 0:
    return func(_ *hugo.Book) error {
      return nil
    }
  case 1:
    return handlers[0]
  default:
    a := handlers[0]
    for _, b := range handlers[1:] {
      a = a.Then(b)
    }
    return a
  }
}

func (a GeneratorHandler) RunOnce(s *bool, b GeneratorHandler) GeneratorHandler {
  return a.Then(func(book *hugo.Book) error {
    if !*s {
      *s = true
      return b(book)
    }
    return nil
  })
}

type Generator struct {
  config     *hugo.Config // Configuration
  bookShelf  *hugo.BookShelf
  generators map[string]GeneratorHandler // Map of available generators
}

func (g *Generator) Name() string {
  return "generator"
}

func (g *Generator) Init(k *kernel.Kernel) error {

  service, err := k.AddService(&hugo.Config{})
  if err != nil {
    return err
  }
  g.config = service.(*hugo.Config)

  service, err = k.AddService(&hugo.BookShelf{})
  if err != nil {
    return err
  }
  g.bookShelf = service.(*hugo.BookShelf)

  return nil
}

func (g *Generator) Start() error {
  g.generators = make(map[string]GeneratorHandler)
  return nil
}

func (g *Generator) Register(n string, h GeneratorHandler) *Generator {
  if _, exists := g.generators[n]; exists {
    panic(fmt.Errorf("GeneratorHandler %s already registered", n))
  }

  g.generators[n] = h
  return g
}

func (g *Generator) Run() error {
  return g.bookShelf.Books().ForEach(
    hugo.WithBook().
        ForEachGenerator(
          hugo.WithBookGenerator().
            Then(g.invokeGenerator)).
        IfExcelPresent(func(book *hugo.Book, excel util.ExcelBuilder) error {
          return excel.FileHandler().
            Write(util.ReferenceFilename(book.ContentPath(), "", "reference.xlsx"), book.Modified())
        }))
}

func (g *Generator) invokeGenerator(book *hugo.Book, n string) error {
  h, exists := g.generators[n]
  if exists {
    return h(book)
  }

  // Log a warning but ignore - could be an invalid config or the generator is not deployed.
  // Originally this was a fatal error, but now we just ignore to allow custom tools to be run
  log.Printf("book %s GeneratorHandler %s is not registered", book.ID, n)
  return nil
}
