package latex

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/latex/parser"
  "github.com/peter-mount/documentation/tools/web"
  "github.com/peter-mount/go-kernel/v2/util/task"
  "log"
  "os"
  "strings"
)

type LaTeX struct {
  bookShelf *hugo.BookShelf `kernel:"inject"`                             // Bookshelf
  enable    *bool           `kernel:"flag,latex,enable LaTeX generation"` // Is LaTeX generation enabled
  worker    task.Queue      `kernel:"worker"`                             // Worker queue
  webserver *web.Webserver  `kernel:"inject"`                             // Webserver
  _         *hugo.Hugo      `kernel:"inject"`                             // access them directly
}

func (l *LaTeX) Start() error {
  if *l.enable {
    return l.bookShelf.Books().ForEach(func(book *hugo.Book) error {
      l.worker.AddPriorityTask(tools.PriorityPDF, func(ctx context.Context) error {
        return l.generate(book)
      })
      return nil
    })
  }

  return nil
}

func appendArg(a []string, flag bool, f string, s ...interface{}) []string {
  if flag {
    return append(a, fmt.Sprintf(f, s...))
  }
  return a
}

func (l *LaTeX) generate(book *hugo.Book) error {
  log.Println("Generating LaTeX for", book.ID)

  url := l.webserver.WebPath("%s/_print/index.html", strings.ToLower(book.WebPath()))

  log.Printf("Retrieving %s", url)
  p, err := parser.Parse(url)
  if err != nil {
    return err
  }

  dest := "public/static/book/" + book.ID + "-tx.tex"
  f, err := os.Create(dest)
  if err != nil {
    return err
  }
  w := NewWriter(f)
  defer w.Close()

  w = w.DocumentClass().
    UsePackage("color").
    UsePackage("framed").
    UsePackage("comment").
    Begin("document")

  for _, n := range p.GetElementByClass("td-main") {
    err = l.parseNode(w, n)
    if err != nil {
      return err
    }
  }

  return nil
}
