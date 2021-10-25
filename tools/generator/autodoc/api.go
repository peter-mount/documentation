package autodoc

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/util/autodoc"
  "github.com/peter-mount/documentation/tools/util/autodoc/asm"
  "github.com/peter-mount/documentation/tools/util/task"
)

// Api manages the collation of API calls in a book
type Api struct {
  m         map[string]*ApiEntry
  api       []*ApiEntry
  extracted bool
}

func NewApi() *Api {
  return &Api{m: make(map[string]*ApiEntry)}
}

func (a *Api) Add(e *ApiEntry) error {
  if ee, exists := a.m[e.Name]; exists {
    return fmt.Errorf("entry %q (%q) already exists with address %q", e.Name, e.Addr, ee.Addr)
  }

  a.m[e.Name] = e
  a.api = append(a.api, e)

  return nil
}

func (a *Api) ForEach(h ApiEntryHandler) error {
  for _, e := range a.api {
    err := h(e)
    if err != nil {
      return err
    }
  }
  return nil
}

// generateSource is a Task for generating source code
func (a *Api) generateSource(ctx context.Context) error {
  book := generator.GetBook(ctx)

  dirName := book.ContentPath("reference/include")
  fileName := "api"

  task.GetQueue(ctx).
      AddTask(task.Of().
        Then(a.SortByAddr).
          Then(autodoc.For(dirName, fileName, book.Modified(), ctx).
            Using(asm.BeebAsm).
            Using(asm.ZAsm).
            InvokeTopic("API", buildHeaderFile).
            Invoke(a.autodocHandler()).
            Do).
        WithContext(ctx, generator.BookKey, autodoc.ResourceManagerKey))

  return nil
}

func (a *Api) autodocHandler() autodoc.Handler {
  return func(ctx context.Context, b autodoc.Builder) error {
    return a.ForEach(func(e *ApiEntry) error {
      b.Function(e.Name, b.Hex(e.Addr), e.Title)
      return nil
    })
  }
}
