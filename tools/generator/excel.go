package generator

import (
  "context"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "path"
  "time"
)

// Excel service that manages multiple Workbooks by ID and ensures they are written
type Excel struct {
  worker   *kernel.Worker // Worker queue
  builders map[string]*provider
}

// provider implements the ExcelProvider interface
type provider struct {
  name     string
  modified time.Time
  builder  util.ExcelBuilder
}

func (p *provider) BuildExcel(f func(builder util.ExcelBuilder) util.ExcelBuilder) error {
  p.builder = f(p.builder)
  return nil
}

func (e *Excel) Name() string {
  return "Excel"
}

func (e *Excel) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&kernel.Worker{})
  if err != nil {
    return err
  }
  e.worker = service.(*kernel.Worker)

  return nil
}

func (e *Excel) Start() error {
  e.builders = make(map[string]*provider)
  return nil
}

// Get returns a named ExcelProvider and ensures it is written to disk.
func (e *Excel) Get(name string, modified time.Time) util.ExcelProvider {
  if provider, exists := e.builders[name]; exists {
    return provider
  }

  provider := &provider{name: name, modified: modified}
  e.builders[name] = provider

  // Add task that will actually write the Excel file
  e.worker.AddPriorityTask(tools.PriorityExcel, provider.task)

  return provider
}

func (p *provider) task(_ context.Context) error {
  if p.builder != nil {
    // NOTE: Use WriteAlways with Excel as we cannot compare an existing version due to xlsx files being zip files
    // so the timestamps inside the zip file are always different causing the generated file to differ
    // even if the actual content is identical.
    return p.builder.
      FileHandler().
      WriteAlways(path.Join("static/static/book/", p.name+".xlsx"), p.modified)
  }
  return nil
}
