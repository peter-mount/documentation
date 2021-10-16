package generator

// Task is a task that the Generator must run once all other Handler's have been run.
// They are usually tasks created by those Handlers.
type Task func() error

type ContextTask func(Context) error

type Context interface {
  AddTask(t Task) Context
  AddPriorityTask(priority int, task Task) Context
  AddContextTask(t ContextTask) Context
  AddContextPriorityTask(priority int, t ContextTask) Context
}

// AddTask appends a Task to be performed once all Handler's have run.
func (g *Generator) AddTask(t Task) Context {
  g.tasks.Add(t)
  return g
}

// AddPriorityTask appends a Task to be performed once all Handler's have run.
func (g *Generator) AddPriorityTask(priority int, task Task) Context {
  g.tasks.AddPriority(priority, task)
  return g
}

func (g *Generator) AddContextTask(t ContextTask) Context {
  return g.AddTask(func() error {
    return t(g)
  })
}

func (g *Generator) AddContextPriorityTask(priority int, t ContextTask) Context {
  return g.AddPriorityTask(priority, func() error {
    return t(g)
  })
}
