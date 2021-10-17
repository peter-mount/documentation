package autodoc

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util/autodoc"
)

// Header is a constant value taken from source, e.g. memory map
type Header struct {
  Label   string // Label of entry
  Value   string // Value of constant
  Comment string // Comment about this entry
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

func (h *Headers) task(book *hugo.Book) generator.Task {
  return func(_ context.Context) error {
    return autodoc.For(book.ContentPath("reference/include"), "headers", book.Modified()).
      Using(autodoc.BeebAsm).
      Using(autodoc.ZAsm).
      InvokeTopic("Headers", book.Autodoc).
      Invoke(h.AutodocHandler()).
      Do()
  }
}

func (h *Headers) AutodocHandler() autodoc.Handler {
  return func(b autodoc.Builder) error {
    return h.ForEach(func(h1 *Header) error {
      if h1.Label == "" && h1.Value == "" && h1.Comment != "" {
        b.Comment(h1.Comment)
      } else {
        b.Header(h1.Label, b.Hex(h1.Value), h1.Comment)
      }
      return nil
    })
  }
}
