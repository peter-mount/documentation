package genlatex

import (
	"context"
	"flag"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/latex"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/go-kernel/v2/log"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"path/filepath"
)

type LaTeX struct {
	Output  *string        `kernel:"flag,o,output filename"`
	handler parser.Handler // Handler to call
}

func (s *LaTeX) Start() error {
	args := flag.Args()
	if len(args) != 1 || *s.Output == "" {
		return fmt.Errorf("syntax: %s -o output url", os.Args[0])
	}
	page := args[0]

	s.handler = latex.New()

	log.Printf("Parsing %s\n", page)

	resp, err := http.Get(page)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ctx := context.Background()
	if *s.Output == "-" {
		ctx = latex.WithContext(os.Stdout, ctx)
	} else {
		log.Println("Creating", *s.Output)
		err = os.MkdirAll(filepath.Dir(*s.Output), 0755)
		if err != nil {
			return err
		}

		f, err := os.Create(*s.Output)
		if err != nil {
			return err
		}
		defer f.Close()
		ctx = latex.WithContext(f, ctx)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	return s.handler.Do(doc, ctx)
}
