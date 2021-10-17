package autodoc

import (
  "context"
  "errors"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util/task"
  "github.com/peter-mount/go-kernel"
  "github.com/peter-mount/go-kernel/util"
)

type Autodoc struct {
  generator *generator.Generator // Generator
  extracted util.Set             // Set of book ID's so that we run once per book
  headers   map[string]*Headers  // Headers per book
}

func (s *Autodoc) Name() string {
  return "Autodoc"
}

func (s *Autodoc) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&generator.Generator{})
  if err != nil {
    return err
  }
  s.generator = service.(*generator.Generator)

  return nil
}

func (s *Autodoc) Start() error {
  s.extracted = util.NewHashSet()
  s.headers = make(map[string]*Headers)

  s.generator.
    Register("autodoc", task.Of(s.extract))

  return nil
}

func (s *Autodoc) getBook(ctx context.Context) (*hugo.Book, error) {
  book, ok := ctx.Value("book").(*hugo.Book)
  if !ok {
    return nil, errors.New("no book present in context")
  }
  return book, nil
}

func (s *Autodoc) getHeaders(ctx context.Context) (*Headers, error) {
  book, err := s.getBook(ctx)
  if err != nil {
    return nil, err
  }

  if h, exists := s.headers[book.ID]; exists {
    return h, nil
  }

  h := NewHeaders()
  s.headers[book.ID] = h

  s.generator.AddTask(h.task(book))

  return h, nil
}
