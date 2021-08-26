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
        Do(ctx)

      if err != nil {
        return err
      }

      *res = buf
      return nil
    }),
  }
}
