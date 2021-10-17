package generator

import (
  "github.com/peter-mount/documentation/tools/hugo"
)

// Handler performs an action against a Book.
// Handler's can either do all of their work in one go or pass Task's to the Generator to be run after
// all other Handler's have run.
type Handler func(*hugo.Book) error

// Then returns a GeneratorHandler that calls this one then if no error the next one forming a chain
func (h Handler) Then(next Handler) Handler {
  return func(b *hugo.Book) error {
    err := h(b)
    if err != nil {
      return err
    }
    return next(b)
  }
}

// HandlerOf returns a Handler that will invoke each supplied GeneratorHandler in a single instance.
// This is a convenience form of calling Then() on each one.
func HandlerOf(handlers ...Handler) Handler {
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

// RunOnce will call a Handler once. It uses a pointer to a boolean to store this state.
// It's useful for simple tasks but should be treated as Deprecated as it only works for one Book not multiple books.
func (h Handler) RunOnce(s *bool, b Handler) Handler {
  return h.Then(func(book *hugo.Book) error {
    if !*s {
      *s = true
      return b(book)
    }
    return nil
  })
}
