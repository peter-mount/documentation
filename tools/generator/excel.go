package generator

import (
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "time"
)

type Excel struct {
  generator *Generator // Generator
  builders  map[string]*provider
}

type provider struct {
  builder util.ExcelBuilder
}

func (p *provider) GetExcel() util.ExcelBuilder {
  return p.builder
}

func (p *provider) SetExcel(b util.ExcelBuilder) {
  p.builder = b
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

func (e *Excel) Get(n string, modified time.Time) util.ExcelProvider {
  if p, exists := e.builders[n]; exists {
    return p
  }

  p := &provider{}
  e.builders[n] = p

  e.generator.AddPriorityTask(500 + len(n), func() error {
    if p.builder != nil {
      return p.builder.
        FileHandler().
        Write("content/"+n+"/reference.xlsx", modified)
    }
    return nil
  })

  return p
}
