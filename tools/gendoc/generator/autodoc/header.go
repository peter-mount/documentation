package autodoc

import (
	"context"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/peter-mount/documentation/tools/gendoc/generator"
	"github.com/peter-mount/documentation/tools/gendoc/util/autodoc"
	"github.com/peter-mount/documentation/tools/gendoc/util/autodoc/asm"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"strings"
)

// Header is a constant value taken from source, e.g. memory map
type Header struct {
	Label   string // Label of entry
	Value   string // Value of constant
	Comment string // Comment about this entry
	Inline  bool   // If true then comment is forced as an inline comment not a new line one
}

type HeaderHandler func(*Header) error

func (h HeaderHandler) Then(b HeaderHandler) HeaderHandler {
	if h == nil {
		return b
	}
	if b == nil {
		return h
	}
	return func(header *Header) error {
		err := h(header)
		if err != nil {
			return err
		}
		return b(header)
	}
}

type Headers struct {
	m map[string]*Header // Map of headers
	a []*Header          // Slice of headers for ordering
}

func NewHeaders() *Headers {
	return &Headers{m: make(map[string]*Header)}
}

var hdrBreak = []string{"<br/>", "<br>", "\n", "\\n"}

func (h *Headers) Add(header *Header) error {
	if header.Label != "" {
		if e, exists := h.m[header.Label]; exists {
			return fmt.Errorf("entry %q already exists with %q ; %q", e.Label, e.Value, e.Comment)
		}

		h.m[header.Label] = header
	}
	h.a = append(h.a, header)

	return nil
}

func (h *Headers) ForEach(f HeaderHandler) error {
	for _, e := range h.a {
		err := f(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// task returns a Task to generate the header include file
func (h *Headers) task(ctx context.Context) error {
	book := generator.GetBook(ctx)

	dirName := book.ContentPath("reference/include")
	fileName := "headers"

	task.GetQueue(ctx).
		AddTask(task.Of().
			Then(autodoc.For(dirName, fileName, book.Modified(), ctx).
				Using(asm.BeebAsm).
				Using(asm.ZAsm).
				InvokeTopic("Headers", buildHeaderFile).
				Invoke(h.AutodocHandler()).
				Do).
			WithContext(ctx, generator.BookKey, autodoc.ResourceManagerKey))

	return nil
}

func (h *Headers) AutodocHandler() autodoc.Handler {
	return func(ctx context.Context, b autodoc.Builder) error {
		return h.ForEach(func(h1 *Header) error {
			comment := stripHtml(splitComment(h1.Comment))
			if !h1.Inline && h1.Label == "" && h1.Value == "" && comment != "" {
				b.Comment(comment)
			} else {
				b.Header(h1.Label, b.Hex(h1.Value), comment)
			}
			return nil
		})
	}
}

// Splits a comment at the first sign of a line break
func splitComment(comment string) string {
	for _, brk := range hdrBreak {
		i := strings.Index(comment, brk)
		if i > -1 {
			return comment[:i]
		}
	}
	return comment
}

// stripHtml removes any html elements that can appear in
func stripHtml(comment string) string {
	if comment == "" {
		return ""
	}

	// Do not allow any html or element content
	return bluemonday.NewPolicy().
		AllowElements().
		SkipElementsContent("sub", "sup").
		Sanitize(comment)
}
