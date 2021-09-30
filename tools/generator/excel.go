package generator

import (
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
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
func (e *Excel) Get(n string, modified time.Time) util.ExcelProvider {
  if p, exists := e.builders[n]; exists {
    return p
  }

  p := &provider{}
  e.builders[n] = p

  e.generator.AddPriorityTask(500+len(n), func() error {
    if p.builder != nil {
      // NOTE: Use WriteAlways with Excel as we cannot compare an existing version due to xlsx files being zip files
      // so the timestamps are always different & the generated file will differ even with identical content.
      return p.builder.
        FileHandler().
        WriteAlways("content/"+n+"/reference.xlsx", modified)
    }
    return nil
  })

  return p
}
