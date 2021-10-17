package task

import (
  "context"
)

// Task is a task that the Generator must run once all other Handler's have been run.
// They are usually tasks created by those Handlers.
type Task func(ctx context.Context) error

// NewTask creates a new Task chain
func NewTask() Task {
  return nil
}

// Then joins two tasks together
func (a Task) Then(b Task) Task {
  if a == nil {
    return b
  }
  if b == nil {
    return a
  }
  return func(ctx context.Context) error {
    if err := a(ctx); err != nil {
      return err
    }
    return b(ctx)
  }
}

// WithValue adds a named value to the context
func (a Task) WithValue(key string, value interface{}) Task {
  return func(ctx context.Context) error {
    return a(context.WithValue(ctx, key, value))
  }
}

func (a Task) Do(ctx context.Context) error {
  if a != nil {
    return a(ctx)
  }
  return nil
}
