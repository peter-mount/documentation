package pdf

import (
  "context"
  "github.com/chromedp/chromedp"
  "log"
)

// Start chrome
func (p *PDF) Start() error {
  // Do nothing
  if *p.baseDir == "" {
    return nil
  }

  log.Println("Starting chromium")

  p.allocCtx, p.allocCancel = chromedp.NewExecAllocator(
    context.Background(),
    chromedp.Headless,
    chromedp.NoSandbox,
    chromedp.NoDefaultBrowserCheck,
    chromedp.NoFirstRun,
    chromedp.DisableGPU,
  )

  p.chromeCtx, p.chromeCancel = chromedp.NewContext(p.allocCtx)

  return nil
}

// Stop chrome
func (p *PDF) Stop() {
  if p.allocCancel != nil {
    p.allocCancel()
  }
  if p.chromeCancel != nil {
    p.chromeCancel()
  }
}

func (p *PDF) run(actions ...chromedp.Action) error {
  return chromedp.Run(p.chromeCtx, actions...)
}
