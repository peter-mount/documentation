package hugo

import (
	"context"
	"github.com/peter-mount/documentation/tools/gendoc"
	"github.com/peter-mount/go-kernel/v2/log"
	util2 "github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/strings"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"os"
	"os/exec"
	"path/filepath"
)

// Hugo runs hugo
type Hugo struct {
	server  *bool      `kernel:"flag,s,Run hugo in server mode"`                       // true to run Hugo in server mode
	cleanup *bool      `kernel:"flag,clean,Cleanup public directories before running"` // true to clean out public directory before running
	draft   *bool      `kernel:"flag,buildDrafts,Build draft pages"`                   // true to build drafts, same as --buildDrafts for hugo
	expired *bool      `kernel:"flag,buildExpired,Build expired pages"`                // true to build expired content, same as --buildExpired for hugo
	future  *bool      `kernel:"flag,buildFuture,Build future pages"`                  // true to build future content, same as --buildFuture for hugo
	worker  task.Queue `kernel:"worker"`                                               // Worker queue
	hugoCmd string     // Path to hugo
}

func (h *Hugo) Start() error {
	if *h.cleanup {
		h.worker.AddPriorityTask(tools.PriorityImmediate, h.cleanupTask)
	}

	h.worker.AddPriorityTask(tools.PriorityHugo, h.run)

	h.hugoCmd = filepath.Join(filepath.Dir(os.Args[0]), "hugo")

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

	stdout := &util2.LogStream{}
	defer stdout.Close()

	cmd := exec.Command(h.hugoCmd, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	return cmd.Run()
}
