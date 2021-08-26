package pdf

import (
  "context"
  "flag"
  "fmt"
  "github.com/chromedp/cdproto/page"
  "github.com/chromedp/chromedp"
  "io/ioutil"
  "log"
)

// Run through args for book id's and generate the PDF's
func (p *PDF) Run() error {
  // Do nothing
  if *p.baseDir == "" {
    return nil
  }

  for _, book := range flag.Args() {
    err := p.generate(book)
    if err != nil {
      return err
    }
  }

  return nil
}

func (p *PDF) generate(book string) error {
  log.Println("Generating PDF for", book)

  url := fmt.Sprintf("http://127.0.0.1:8080/%s/_print/", book)

  // capture pdf
  var buf []byte
  err := p.run(p.printToPDF(url, &buf))
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(book+".pdf", buf, 0644)
  if err != nil {
    return err
  }

  return nil
}

// print a specific pdf page.
func (p *PDF) printToPDF(url string, res *[]byte) chromedp.Tasks {
  return chromedp.Tasks{
    chromedp.Navigate(url),
    chromedp.ActionFunc(func(ctx context.Context) error {
      buf, _, err := page.PrintToPDF().
        WithPrintBackground(false).
        WithMarginTop(0.5).
        WithMarginBottom(0.5).
        //WithMarginLeft(.1).
        //WithMarginRight(.1).
        WithLandscape(false).
        WithDisplayHeaderFooter(true).
        WithHeaderTemplate("<div style='font-size: 5px;text-align:center;width:100%;margin-left:3em;margin-right:3em;'><span class='title'></span></div>").
        WithFooterTemplate("<div style='font-size: 5px;width:100%;margin-left:3em;margin-right:3em;'><span class='date'></span> <span style='float:right'><span class='pageNumber'></span>&nbsp;of&nbsp;<span class='totalPages'></span></span></div>").
        Do(ctx)

      if err != nil {
        return err
      }

      *res = buf
      return nil
    }),
  }
}
