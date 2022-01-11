package hugo

import (
  "context"
  "flag"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "log"
  "os"
  "os/exec"
)

// Hugo runs hugo
type Hugo struct {
  server  *bool          // true to run Hugo in server mode
  cleanup *bool          // true to clean out public directory before running
  worker  *kernel.Worker // Worker queue
}

func (h *Hugo) Name() string {
  return "hugo"
}

func (h *Hugo) Init(k *kernel.Kernel) error {
  h.server = flag.Bool("s", false, "Run hugo in server mode")
  h.cleanup = flag.Bool("clean", false, "Cleanup public directory before running")

  service, err := k.AddService(&kernel.Worker{})
  if err != nil {
    return err
  }
  h.worker = service.(*kernel.Worker)

  return k.DependsOn(&PostCSS{})
}

func (h *Hugo) Start() error {
  if *h.cleanup {
    h.worker.AddPriorityTask(tools.PriorityImmediate, h.cleanupTask)
  }

  h.worker.AddPriorityTask(tools.PriorityHugo, h.run)

  return nil
}

// Remove all of our temp dirs
func (h *Hugo) cleanupTask(_ context.Context) error {
  log.Println("Cleaning up directories")

  return util.StringSliceOf(
    "public",
    "static/static/book",
    "static/static/chipref",
    "static/static/gen",
  ).ForEach(os.RemoveAll)
}

func (h *Hugo) run(_ context.Context) error {
  log.Println("Running hugo")

  var args []string
  if *h.server {
    args = append(args, "server", "--bind=0.0.0.0", "--port=1313")
  }

  stdout := &util.LogStream{}
  defer stdout.Close()

  cmd := exec.Command("hugo", args...)
  cmd.Stdout = stdout
  cmd.Stderr = stdout

  return cmd.Run()
}
