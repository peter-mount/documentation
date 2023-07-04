package svg

import (
	"context"
	"github.com/peter-mount/documentation/tools/gendoc"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type SVG struct {
	worker task.Queue `kernel:"worker"` // Worker queue
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
