package svg

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/util/walk"
  "github.com/peter-mount/go-kernel"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
)

type SVG struct {
  generator *generator.Generator // Generator
}

func (s *SVG) Name() string {
  return "SVG"
}

func (s *SVG) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&generator.Generator{})
  if err != nil {
    return err
  }
  s.generator = service.(*generator.Generator)

  return nil
}

func (s *SVG) Start() error {
  return walk.NewPathWalker().
    Then(s.processFile).
    PathHasSuffix(".yaml").
    IsFile().
    Walk("config")
}

func (s *SVG) processFile(path string, info os.FileInfo) error {
  s.generator.AddPriorityTask(70, func(ctx context.Context) error {
    buf, err := ioutil.ReadFile(path)
    if err != nil {
      return err
    }

    file := &File{}
    err = yaml.Unmarshal(buf, file)
    if err != nil {
      return err
    }

    return file.FileHandler().
      Write(file.FileName, info.ModTime())
  })

  return nil
}
