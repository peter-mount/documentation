package autodoc

import (
	"context"
	"github.com/peter-mount/documentation/tools/gendoc"
	"github.com/peter-mount/documentation/tools/gendoc/generator"
	"github.com/peter-mount/documentation/tools/gendoc/util/autodoc"
	"github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/task"
)

type Autodoc struct {
	generator       *generator.Generator     `kernel:"inject"` // Generator
	resourceManager *autodoc.ResourceManager `kernel:"inject"` // ResourceManager
	extracted       util.Set                 // Set of book ID's so that we run once per book
	headers         map[string]*Headers      // Headers per book
	apis            map[string]*Api          // Apis per book
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
		AddPriorityTask(tools.PriorityApi, task.Of(h.task).
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
		AddPriorityTask(tools.PriorityApi, task.Of().
			Then(a.generateResource).
			Then(a.generateSource).
			Then(a.generateIndex).
			WithValue(generator.BookKey, book).
			WithValue(autodoc.ResourceManagerKey, s.resourceManager))

	return a
}
