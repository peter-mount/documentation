package pdf

import (
  "context"
  "github.com/chromedp/cdproto/page"
  "github.com/chromedp/chromedp"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/go-kernel"
  "log"
)

// PDF tool that handles the generation of PDF documentation of a "book"
type PDF struct {
  config   *hugo.Config   // Config
  chromium *hugo.Chromium // Chromium browser
}

func (p *PDF) Name() string {
  return "PDF"
}

func (p *PDF) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&hugo.Config{})
  if err != nil {
    return err
  }
  p.config = service.(*hugo.Config)

  service, err = k.AddService(&hugo.Chromium{})
  if err != nil {
    return err
  }
  p.chromium = service.(*hugo.Chromium)

  // These we need these services to be running before us
  return k.DependsOn(&hugo.Generator{}, &hugo.Webserver{}, &hugo.Hugo{})
}

// Run through args for book id's and generate the PDF's
func (p *PDF) Run() error {

  for _, book := range p.config.Books {
    err := p.generate(book)
    if err != nil {
      return err
    }
  }

  return nil
}

func (p *PDF) generate(book *hugo.Book) error {
  log.Println("Generating PDF for", book.ID)

  // capture pdf
  var buf []byte
  err := p.chromium.Run(p.printToPDF(book, &buf))
  if err != nil {
    return err
  }

  return util.WriteFile(book.ID+".pdf", buf, book.Modified())
}

// print a specific pdf page.
func (p *PDF) printToPDF(book *hugo.Book, res *[]byte) chromedp.Tasks {
  url := p.config.Webpath("%s/_print/", book.ID)

  pdf := book.PDF

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
        WithDisplayHeaderFooter(!pdf.DisableHeaderFooter).
        WithHeaderTemplate(book.Expand(pdf.Header)).
        WithFooterTemplate(book.Expand(pdf.Footer)).
        Do(ctx)

      if err != nil {
        return err
      }

      *res = buf
      return nil
    }),
  }
}
