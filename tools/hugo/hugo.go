package hugo

import (
  "flag"
  "github.com/peter-mount/go-kernel"
  "log"
  "os/exec"
)

// Hugo runs hugo & provides a webserver if required for other tasks
type Hugo struct {
  hugo *bool // Run Hugo
}

func (h *Hugo) Name() string {
  return "hugo"
}

func (h *Hugo) Init(k *kernel.Kernel) error {
  h.hugo = flag.Bool("hugo", false, "Run Hugo")

  return nil
}

func (h *Hugo) Run() error {
  // Do nothing
  if !*h.hugo {
    return nil
  }

  log.Println("Running hugo")
  stdout := &LogStream{}
  defer stdout.Close()

  cmd := exec.Command("hugo")
  cmd.Stdout = stdout
  cmd.Stderr = stdout
  err := cmd.Run()

  return err
}
