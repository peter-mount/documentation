package task

import (
  "context"
)

// Task is a task that the Generator must run once all other Handler's have been run.
// They are usually tasks created by those Handlers.
type Task func(ctx context.Context) error

// Of creates a new Task forming a chain of the provided tasks
func Of(tasks ...Task) Task {
  var task Task
  for _, b := range tasks {
    task = task.Then(b)
  }
  return task
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
    err := a(ctx)
    if err == nil {
      err = b(ctx)
    }
    return err
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

// RunOnce will invoke a task exactly once.
// It uses a pointer to a boolean to store this state.
// It's useful for simple tasks but should be treated as Deprecated.
// Currently, here as Book still uses it as it only works for one Book not multiple books.
func (a Task) RunOnce(f *bool, t Task) Task {
  return a.Then(func(ctx context.Context) error {
    if !*f {
      *f = true
      return t(ctx)
    }
    return nil
  })
}
