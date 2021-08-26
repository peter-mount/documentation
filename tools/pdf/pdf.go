package pdf

import (
  "context"
  "flag"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/go-kernel"
)

// PDF tool that handles the generation of PDF documentation of a "book"
type PDF struct {
  webserver    *hugo.Webserver    // We need the webserver
  baseDir      *string            // Directory to place pdf
  chromeCtx    context.Context    // Chrome context
  chromeCancel context.CancelFunc // close chrome
  allocCtx     context.Context    // alloc context
  allocCancel  context.CancelFunc // alloc close
}

func (p *PDF) Name() string {
  return "PDF"
}

func (p *PDF) Init(k *kernel.Kernel) error {
  p.baseDir = flag.String("pdf", "", "Directory for PDF generation")

  service, err := k.AddService(&hugo.Webserver{})
  if err != nil {
    return err
  }
  p.webserver = service.(*hugo.Webserver)

  return nil
}

func (p *PDF) PostInit() error {
  // If we are to do anything then enable the webserver
  p.webserver.Enable(*p.baseDir != "")

  return nil
}
