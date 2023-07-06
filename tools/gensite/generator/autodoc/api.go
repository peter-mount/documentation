package autodoc

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/gensite/generator"
	"github.com/peter-mount/documentation/tools/gensite/util"
	"github.com/peter-mount/documentation/tools/gensite/util/autodoc"
	"github.com/peter-mount/documentation/tools/gensite/util/autodoc/asm"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"os"
	"path"
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
	if len(a.api) > 0 {
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

	}

	return nil
}

func (a *Api) generateResource(ctx context.Context) error {
	book := generator.GetBook(ctx)

	dirName := book.ContentPath("reference")

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return err
	}

	task.GetQueue(ctx).
		AddTask(
			autodoc.GenerateReferenceIndices(path.Join(dirName, "_index.html"), book.Modified()).
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

// generateIndex generates the API indices
func (a *Api) generateIndex(ctx context.Context) error {
	if len(a.api) > 0 {

		task.GetQueue(ctx).
			AddTask(task.Of().
				Then(a.SortByAddr).
				Then(a.generateIndexFile).
				WithValue("filename", "api").
				WithValue("title", "API by Address").
				WithContext(ctx, generator.BookKey, autodoc.ResourceManagerKey)).
			AddTask(task.Of().
				Then(a.SortByName).
				Then(a.generateIndexFile).
				WithValue("filename", "apiname").
				WithValue("title", "API by Name").
				WithContext(ctx, generator.BookKey, autodoc.ResourceManagerKey))

	}
	return nil
}

func (a *Api) generateIndexFile(ctx context.Context) error {
	book := generator.GetBook(ctx)

	//dirName := book.ContentPath("reference")
	fileName := ctx.Value("filename").(string)
	title := ctx.Value("title").(string)

	r := Output{Nometa: true}
	for _, o := range a.api {
		r.Api = append(r.Api, o.params)
	}

	return util.ReferenceFileBuilder(title, title, "manual", 10, book.Modified()).
		Yaml(r).
		WrapAsFrontMatter().
		FileHandler().
		Write(util.ReferenceFilename(book.ContentPath(), fileName, "_index.html"), book.Modified())
}
