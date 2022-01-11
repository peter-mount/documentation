package autodoc

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/util/autodoc"
  "github.com/peter-mount/go-kernel"
  "github.com/peter-mount/go-kernel/util"
  "github.com/peter-mount/go-kernel/util/task"
)

type Autodoc struct {
  generator       *generator.Generator     // Generator
  resourceManager *autodoc.ResourceManager // ResourceManager
  extracted       util.Set                 // Set of book ID's so that we run once per book
  headers         map[string]*Headers      // Headers per book
  apis            map[string]*Api          // Apis per book
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

  service, err = k.AddService(&autodoc.ResourceManager{})
  if err != nil {
    return err
  }
  s.resourceManager = service.(*autodoc.ResourceManager)

  return nil
}

func (s *Autodoc) Start() error {
  s.extracted = util.NewHashSet()
  s.headers = make(map[string]*Headers)
  s.apis = make(map[string]*Api)

  s.generator.
      Register("autodoc", task.Of(s.extract).
        WithValue(autodoc.ResourceManagerKey, s.resourceManager))

  return nil
}

func (s *Autodoc) getHeaders(ctx context.Context) *Headers {
  book := generator.GetBook(ctx)

  if h, exists := s.headers[book.ID]; exists {
    return h
  }

  h := NewHeaders()
  s.headers[book.ID] = h

  task.GetQueue(ctx).
      AddPriorityTask(20, task.Of(h.task).
        WithValue(generator.BookKey, book).
        WithValue(autodoc.ResourceManagerKey, s.resourceManager))

  return h
}

func (s *Autodoc) GetApi(ctx context.Context) *Api {
  book := generator.GetBook(ctx)

  if a, exists := s.apis[book.ID]; exists {
    return a
  }

  a := NewApi()
  s.apis[book.ID] = a

  task.GetQueue(ctx).
      AddPriorityTask(20, task.Of().
        Then(a.generateResource).
        Then(a.generateSource).
        Then(a.generateIndex).
        WithValue(generator.BookKey, book).
        WithValue(autodoc.ResourceManagerKey, s.resourceManager))

  return a
}
