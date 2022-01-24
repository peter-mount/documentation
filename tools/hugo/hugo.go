package hugo

import (
  "context"
  "flag"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/strings"
  "github.com/peter-mount/go-kernel"
  "log"
  "os"
  "os/exec"
)

// Hugo runs hugo
type Hugo struct {
  server  *bool          // true to run Hugo in server mode
  cleanup *bool          // true to clean out public directory before running
  draft   *bool          // true to build drafts, same as --buildDrafts for hugo
  expired *bool          // true to build expired content, same as --buildExpired for hugo
  future  *bool          // true to build future content, same as --buildFuture for hugo
  worker  *kernel.Worker // Worker queue
}

func (h *Hugo) Name() string {
  return "hugo"
}

func (h *Hugo) Init(k *kernel.Kernel) error {
  h.server = flag.Bool("s", false, "Run hugo in server mode")
  h.cleanup = flag.Bool("clean", false, "Cleanup public directory before running")
  h.draft = flag.Bool("buildDrafts", false, "Build draft pages")
  h.expired = flag.Bool("buildExpired", false, "Build expired pages")
  h.future = flag.Bool("buildFuture", false, "Build future pages")

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

  return strings.Of(
    "public",
    "static/static/book",
    "static/static/chipref",
    "static/static/gen",
  ).ForEach(os.RemoveAll)
}

func appendArg(a []string, flag bool, s ...string) []string {
  if flag {
    return append(a, s...)
  }
  return a
}

func (h *Hugo) run(_ context.Context) error {
  log.Println("Running hugo")

  var args []string
  args = appendArg(args, *h.server, "server", "--bind=0.0.0.0", "--port=1313")
  args = appendArg(args, *h.draft, "--buildDrafts")
  args = appendArg(args, *h.expired, "--buildExpired")
  args = appendArg(args, *h.future, "--buildFuture")

  stdout := &util.LogStream{}
  defer stdout.Close()

  cmd := exec.Command("hugo", args...)
  cmd.Stdout = stdout
  cmd.Stderr = stdout

  return cmd.Run()
}
