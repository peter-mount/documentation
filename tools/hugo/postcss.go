package hugo

import (
  "fmt"
  util2 "github.com/peter-mount/go-kernel/v2/util"
  "log"
  "os"
  "os/exec"
)

// PostCSS ensures that the postcss plugins are installed
// After NodeJS 15 they cannot be globally installed, so we need them
// at the same directory as the repository.
//
// As we can't mount both at the same time we cannot do this as part of
// the container build, so we do it here.
//
// TODO find a way to do this internally in the build
type PostCSS struct {
  installed bool
}

const nodeModules = "node_modules"

// Start checks nodeModules does not exist and if it doesn't installs the
// modules
func (p *PostCSS) Start() error {
  if fi, err := os.Stat(nodeModules); os.IsNotExist(err) {
    return p.install()
  } else if err != nil {
    return err
  } else if !fi.IsDir() {
    return fmt.Errorf("%s is not a directory", nodeModules)
  }
  return nil
}

func (p *PostCSS) Stop() {
  if p.installed {
    err := p.uninstall()
    if err != nil {
      log.Println(err)
    }
  }
}

func (p *PostCSS) install() error {
  log.Println("Installing postcss plugins")

  stdout := &util2.LogStream{}
  defer stdout.Close()

  //cmd := exec.Command("npm", "install", "autoprefixer")
  cmd := exec.Command("cp", "-pr", "/root/"+nodeModules, nodeModules)
  cmd.Stdout = stdout
  cmd.Stderr = stdout

  err := cmd.Run()
  if err == nil {
    p.installed = true
  }
  return err
}

func (p *PostCSS) uninstall() error {
  log.Println("Uninstalling postcss plugins")

  stdout := &util2.LogStream{}
  defer stdout.Close()

  cmd := exec.Command("rm", "-rf", nodeModules)
  cmd.Stdout = stdout
  cmd.Stderr = stdout

  return cmd.Run()
}
