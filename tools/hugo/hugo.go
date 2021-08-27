package hugo

import (
  "github.com/peter-mount/go-kernel"
  "log"
  "os/exec"
)

// Hugo runs hugo & provides a webserver if required for other tasks
type Hugo struct {
}

func (h *Hugo) Name() string {
  return "hugo"
}

func (h *Hugo) Init(k *kernel.Kernel) error {
  // Depend on Generator running first
  return k.DependsOn(&Generator{})
}

func (h *Hugo) Run() error {
  log.Println("Running hugo")
  stdout := &LogStream{}
  defer stdout.Close()

  cmd := exec.Command("hugo")
  cmd.Stdout = stdout
  cmd.Stderr = stdout
  err := cmd.Run()

  return err
}
