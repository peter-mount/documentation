package generator

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/go-kernel"
  "log"
)

// GeneratorHandler performs an action against a Book
type GeneratorHandler func(*hugo.Book) error

// Then returns a GeneratorHandler that calls this one then if no error the next one forming a chain
func (h GeneratorHandler) Then(next GeneratorHandler) GeneratorHandler {
  return func(b *hugo.Book) error {
    err := h(b)
    if err == nil {
      err = next(b)
    }
    return err
  }
}

// NopGeneratorHandler is a GeneratorHandler that does nothing
func NopGeneratorHandler(_ *hugo.Book) error {
  return nil
}

// GeneratorHandlerChain returns a GeneratorHandler that will invoke each supplied GeneratorHandler in a single instance.
// This is a convenience form of calling Then() on each one.
func GeneratorHandlerChain(handlers ...GeneratorHandler) GeneratorHandler {
  if len(handlers) == 0 {
    return NopGeneratorHandler
  }

  a := handlers[0]
  for _, b := range handlers[1:] {
    a = a.Then(b)
  }
  return a
}

type Generator struct {
  config     *hugo.Config                // Configuration
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
  return g.config.Books.ForEach(func(book *hugo.Book) error {
    return book.Generate.ForEach(func(n string) error {
      h, exists := g.generators[n]
      if exists {
        return h(book)
      }

      // Log a warning but ignore - could be an invalid config or the generator is not deployed.
      // Originally this was a fatal error, but now we just ignore to allow custom tools to be run
      log.Printf("book %s GeneratorHandler %s is not registered", book.ID, n)
      return nil
    })
  })
}
