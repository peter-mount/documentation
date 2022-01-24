package pdf

import (
  "context"
  "github.com/chromedp/cdproto/page"
  "github.com/chromedp/chromedp"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/web"
  "github.com/peter-mount/go-kernel/util/task"
  "log"
  "strings"
)

// PDF tool that handles the generation of PDF documentation of a "book"
type PDF struct {
  config    *hugo.Config    `kernel:"inject"`                        // Config
  bookShelf *hugo.BookShelf `kernel:"inject"`                        // Bookshelf
  chromium  *web.Chromium   `kernel:"inject"`                        // Chromium browser
  enable    *bool           `kernel:"flag,p,disable PDF generation"` // Is PDF generation enabled
  worker    task.Queue      `kernel:"worker"`                        // Worker queue
  _         *web.Webserver  `kernel:"inject"`                        // we need these to deploy before this but don't
  _         *hugo.Hugo      `kernel:"inject"`                        // access them directly
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
