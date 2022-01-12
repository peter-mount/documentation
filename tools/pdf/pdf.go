package pdf

import (
  "context"
  "flag"
  "github.com/chromedp/cdproto/page"
  "github.com/chromedp/chromedp"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/web"
  "github.com/peter-mount/go-kernel"
  "log"
  "strings"
)

// PDF tool that handles the generation of PDF documentation of a "book"
type PDF struct {
  config    *hugo.Config // Config
  bookShelf *hugo.BookShelf
  chromium  *web.Chromium  // Chromium browser
  enable    *bool          // Is PDF generation enabled
  worker    *kernel.Worker // Worker queue
}

func (p *PDF) Name() string {
  return "PDF"
}

func (p *PDF) Init(k *kernel.Kernel) error {
  p.enable = flag.Bool("p", false, "disable pdf generation")

  service, err := k.AddService(&hugo.Config{})
  if err != nil {
    return err
  }
  p.config = service.(*hugo.Config)

  service, err = k.AddService(&web.Chromium{})
  if err != nil {
    return err
  }
  p.chromium = service.(*web.Chromium)

  service, err = k.AddService(&hugo.BookShelf{})
  if err != nil {
    return err
  }
  p.bookShelf = service.(*hugo.BookShelf)

  service, err = k.AddService(&kernel.Worker{})
  if err != nil {
    return err
  }
  p.worker = service.(*kernel.Worker)

  // We need a webserver & must run after hugo
  return k.DependsOn(&web.Webserver{}, &hugo.Hugo{})
}

func (p *PDF) Start() error {
  if *p.enable {
    return nil
  }

  return p.bookShelf.Books().ForEach(func(book *hugo.Book) error {
    p.worker.AddPriorityTask(tools.PriorityPDF, func(ctx context.Context) error {
      log.Println("Generating PDF for", book.ID)
      return p.chromium.Run(p.printToPDF(book))
    })
    return nil
  })
}

// print a specific pdf page.
func (p *PDF) printToPDF(book *hugo.Book) chromedp.Tasks {
  url := p.config.WebPath("%s/_print/", strings.ToLower(book.WebPath()))

  pdf := p.config.PDF

  return chromedp.Tasks{
    chromedp.Navigate(url),
    chromedp.ActionFunc(func(ctx context.Context) error {
      buf, _, err := page.PrintToPDF().
        WithPrintBackground(pdf.PrintBackground).
        WithMarginTop(pdf.Margin.Top).
        WithMarginBottom(pdf.Margin.Bottom).
        WithMarginLeft(pdf.Margin.Left).
        WithMarginRight(pdf.Margin.Right).
        WithLandscape(pdf.Landscape).
        WithPaperWidth(pdf.Width).
        WithPaperHeight(pdf.Height).
        WithDisplayHeaderFooter(!pdf.DisableHeaderFooter).
        WithHeaderTemplate(book.Expand(pdf.Header)).
        WithFooterTemplate(book.Expand(pdf.Footer)).
        Do(ctx)

      if err != nil {
        return err
      }

      return util.ByteFileHandler(buf).
        Write("public/static/book/"+book.ID+".pdf", book.Modified())
    }),
  }
}
