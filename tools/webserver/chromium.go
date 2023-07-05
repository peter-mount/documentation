package webserver

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/peter-mount/go-kernel/v2/log"
)

// Chromium embeds a headless chromium browser
type Chromium struct {
	chromeCtx    context.Context    // Chrome context
	chromeCancel context.CancelFunc // close chrome
	allocCtx     context.Context    // alloc context
	allocCancel  context.CancelFunc // alloc close
}

func (c *Chromium) Start() error {
	log.Println("Starting chromium")

	c.allocCtx, c.allocCancel = chromedp.NewExecAllocator(
		context.Background(),
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoFirstRun,
		chromedp.DisableGPU,
	)

	c.chromeCtx, c.chromeCancel = chromedp.NewContext(c.allocCtx)

	return nil
}

// Stop chrome
func (c *Chromium) Stop() {
	log.Println("Stopping chromium")
	c.allocCancel()
	c.chromeCancel()
}

func (c *Chromium) Run(actions ...chromedp.Action) error {
	return chromedp.Run(c.chromeCtx, actions...)
}
