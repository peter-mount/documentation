package pdf

import (
  "context"
  "github.com/chromedp/cdproto/page"
  "github.com/chromedp/chromedp"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/go-kernel"
  "io/ioutil"
  "log"
  "os"
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

  // Just depend on Webserver
  _, err = k.AddService(&hugo.Webserver{})
  return err
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

  fileName := book.ID + ".pdf"
  fileTime := book.Modified()

  err = ioutil.WriteFile(fileName, buf, 0644)
  if err != nil {
    return err
  }

  err = os.Chtimes(fileName, fileTime, fileTime)
  if err != nil {
    return err
  }

  return nil
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
