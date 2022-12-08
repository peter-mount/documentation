package latex

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/latex/parser"
  "github.com/peter-mount/documentation/tools/web"
  util2 "github.com/peter-mount/go-kernel/v2/util"
  "github.com/peter-mount/go-kernel/v2/util/task"
  "log"
  "os"
  "os/exec"
  "path/filepath"
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
      if book.ID == "test" {
        l.worker.AddPriorityTask(tools.PriorityPDF, func(ctx context.Context) error {
          return l.generate(book)
        })
      }
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
  bookDir := filepath.Join("public/static/book", book.ID)
  log.Println("Creating", bookDir)
  err := os.MkdirAll(bookDir, 0777)
  if err != nil {
    return err
  }

  tex := filepath.Join(bookDir, book.ID+".tex")
  //pdf := filepath.Join(bookDir, book.ID+".pdf")

  err = l.generateTeX(book, tex)
  if err == nil {
    err = l.generatePdf(tex)
  }
  return err
}

func (l *LaTeX) generateTeX(book *hugo.Book, tex string) error {
  log.Println("Generating LaTeX for", book.ID)

  url := l.webserver.WebPath("%s/_print/index.html", strings.ToLower(book.WebPath()))

  log.Printf("Retrieving %s", url)
  p, err := parser.Parse(url)
  if err != nil {
    return err
  }

  log.Println("Creating", tex)
  f, err := os.Create(tex)
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

func (l *LaTeX) generatePdf(tex string) error {
  tmpDir := filepath.Dir(tex)

  stdout := &util2.LogStream{}
  defer stdout.Close()

  cmd := exec.Command("latex", "-interaction=nonstopmode", "-output-directory="+tmpDir, "-output-format=pdf", tex)
  cmd.Stdout = stdout
  cmd.Stderr = stdout

  return cmd.Run()
}
