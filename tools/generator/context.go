package generator

import "context"

// Task is a task that the Generator must run once all other Handler's have been run.
// They are usually tasks created by those Handlers.
type Task func(context.Context) error

type TaskContext interface {
  AddTask(t Task) TaskContext
  AddPriorityTask(priority int, task Task) TaskContext
}

const (
  ctxKey = "GeneratorTaskContext"
)

// GetTaskContext returns the TaskContext contained in the passed Context
func GetTaskContext(ctx context.Context) TaskContext {
  if tc, ok := ctx.Value(ctxKey).(TaskContext); ok {
    return tc
  }

  return nil
}

// AddTask appends a Task to be performed once all Handler's have run.
func (g *Generator) AddTask(t Task) TaskContext {
  g.tasks.Add(t)
  return g
}

// AddPriorityTask appends a Task to be performed once all Handler's have run.
func (g *Generator) AddPriorityTask(priority int, task Task) TaskContext {
  g.tasks.AddPriority(priority, task)
  return g
}
