package generator

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util/task"
  "github.com/peter-mount/go-kernel"
  "log"
)

// Generator is a kernel Service which handles the generation of content based on page metadata.
type Generator struct {
  config     *hugo.Config // Configuration
  bookShelf  *hugo.BookShelf
  generators map[string]Handler // Map of available generators
  tasks      task.Queue
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
  g.generators = make(map[string]Handler)
  g.tasks = task.NewQueue()
  return nil
}

// Register a named Handler
func (g *Generator) Register(n string, h Handler) *Generator {
  if _, exists := g.generators[n]; exists {
    panic(fmt.Errorf("GeneratorHandler %s already registered", n))
  }

  g.generators[n] = h
  return g
}

func (g *Generator) Run() error {
  if err := g.bookShelf.Books().ForEach(context.Background(), g.invokeBook); err != nil {
    return err
  }

  return task.Run(g.tasks, context.Background())
}

func (g *Generator) AddTask(task task.Task) task.Queue {
  return g.tasks.AddTask(task)
}

func (g *Generator) AddPriorityTask(priority int, task task.Task) task.Queue {
  return g.tasks.AddPriorityTask(priority, task)
}

func (g *Generator) invokeBook(_ context.Context, book *hugo.Book) error {
  g.AddTask(func(ctx context.Context) error {
    return hugo.WithBook().
        ForEachGenerator(
          hugo.WithBookGenerator().
            Then(g.invokeGenerator)).
      Do(ctx, book)
  })

  return nil
}

func (g *Generator) invokeGenerator(ctx context.Context, book *hugo.Book, n string) error {
  h, exists := g.generators[n]
  if exists {
    return h(book)
  }

  // Log a warning but ignore - could be an invalid config or the generator is not deployed.
  // Originally this was a fatal error, but now we just ignore to allow custom tools to be run
  log.Printf("book %s GeneratorHandler %s is not registered", book.ID, n)
  return nil
}
