package svg

import (
  "context"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/util/walk"
  "github.com/peter-mount/go-kernel"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
)

type SVG struct {
  worker *kernel.Worker // Worker queue
}

func (s *SVG) Name() string {
  return "SVG"
}

func (s *SVG) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&kernel.Worker{})
  if err != nil {
    return err
  }
  s.worker = service.(*kernel.Worker)

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
  s.worker.AddPriorityTask(tools.PrioritySVG, func(ctx context.Context) error {
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
