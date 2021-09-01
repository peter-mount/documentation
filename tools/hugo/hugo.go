package hugo

import (
  "flag"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "log"
  "os/exec"
)

// Hugo runs hugo
type Hugo struct {
  server *bool // true to run Hugo in server mode
}

func (h *Hugo) Name() string {
  return "hugo"
}

func (h *Hugo) Init(_ *kernel.Kernel) error {
  h.server = flag.Bool("s", false, "Run hugo in server mode")
  return nil
}

func (h *Hugo) Run() error {
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
