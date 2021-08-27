package hugo

import (
  "github.com/peter-mount/documentation/tools/util"
  "log"
  "os/exec"
)

// Hugo runs hugo
type Hugo struct {
}

func (h *Hugo) Name() string {
  return "hugo"
}

func (h *Hugo) Run() error {
  log.Println("Running hugo")

  stdout := &util.LogStream{}
  defer stdout.Close()

  cmd := exec.Command("hugo")
  cmd.Stdout = stdout
  cmd.Stderr = stdout

  return cmd.Run()
}
