package generator

import (
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "path"
  "time"
)

// Excel service that manages multiple Workbooks by ID and ensures they are written
type Excel struct {
  generator *Generator // Generator
  builders  map[string]*provider
}

// provider implements the ExcelProvider interface
type provider struct {
  builder util.ExcelBuilder
}

func (p *provider) BuildExcel(f func(builder util.ExcelBuilder) util.ExcelBuilder) error {
  p.builder = f(p.builder)
  return nil
}

func (e *Excel) Name() string {
  return "Excel"
}

func (e *Excel) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&Generator{})
  if err != nil {
    return err
  }
  e.generator = service.(*Generator)

  return nil
}

func (e *Excel) Start() error {
  e.builders = make(map[string]*provider)
  return nil
}

// Get returns a named ExcelProvider and ensures its written to disk.
func (e *Excel) Get(name string, modified time.Time) util.ExcelProvider {
  if provider, exists := e.builders[name]; exists {
    return provider
  }

  provider := &provider{}
  e.builders[name] = provider

  e.generator.AddPriorityTask(500, func() error {
    if provider.builder != nil {
      // NOTE: Use WriteAlways with Excel as we cannot compare an existing version due to xlsx files being zip files
      // so the timestamps inside the zip file are always different causing the generated file to differ
      // even if the actual content is identical.
      return provider.builder.
        FileHandler().
        WriteAlways(path.Join("static/static/book/", name+".xlsx"), modified)
    }
    return nil
  })

  return provider
}
