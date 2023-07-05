package latex

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/gensite"
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/documentation/tools/gensite/latex/css"
	"github.com/peter-mount/documentation/tools/gensite/latex/parser"
	"github.com/peter-mount/documentation/tools/gensite/latex/util"
	"github.com/peter-mount/documentation/tools/webserver"
	"github.com/peter-mount/go-kernel/v2/log"
	util2 "github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type LaTeX struct {
	bookShelf *hugo.BookShelf      `kernel:"inject"` // Bookshelf
	css       *css.Css             `kernel:"inject"`
	enable    *bool                `kernel:"flag,latex,enable LaTeX generation"` // Is LaTeX generation enabled
	worker    task.Queue           `kernel:"worker"`                             // Worker queue
	webserver *webserver.Webserver `kernel:"inject"`                             // Webserver
	_         *hugo.Hugo           `kernel:"inject"`                             // access them directly
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
	w := util.NewWriter(f)
	defer w.Close()

	w.WriteString("%%      Book: %s\n%%     Title: %s\n%% Generated: %s\n",
		book.ID,
		book.Title,
		time.Now().UTC().Format(time.RFC3339),
	)

	w.WriteString("\\documentclass[a4paper,10pt,reqno,twoside,onecolumn,final,openright]{book}\n").
		UsePackage("amsmath,amssymb,amsfonts"). // Typical math packages
		UsePackage("graphics").                 // Allow graphics
		UsePackage("color").                    // Coloured text & backgrounds
		//UsePackage("framed").
		//UsePackage("comment").
		UsePackage("tabularx"). // Fixed width tables
		UsePackage("multirow"). // Needed for rowspan
		WriteString("\\DeclareGraphicsExtensions{.pdf}\n").
		WriteString("\\parindent 0cm\n").
		WriteString("\\parskip 0.2cm\n").
		WriteString("\\topmargin 0.2cm\n").
		WriteString("\\oddsidemargin 1cm\n").
		WriteString("\\evensidemargin 0.5cm\n").
		WriteString("\\textwidth 19cm\n").
		WriteString("\\textheight 28cm\n").
		WriteString("\\begin{document}\n")

	for _, n := range p.GetElementByClass("td-main") {
		err = l.parseNode(w, n)
		if err != nil {
			return err
		}
	}

	w.WriteString("\\end{document}")
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
