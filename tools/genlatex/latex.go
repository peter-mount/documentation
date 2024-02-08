package genlatex

import (
	"context"
	"flag"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/latex"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/go-build/version"
	"github.com/peter-mount/go-kernel/v2/log"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"path/filepath"
)

type LaTeX struct {
	Output     *string `kernel:"flag,o,output filename"`
	Stylesheet *string `kernel:"flag,style,Stylesheet"`
}

func (s *LaTeX) Start() error {
	log.Println(version.Version)

	args := flag.Args()
	if len(args) != 1 || *s.Output == "" {
		return fmt.Errorf("syntax: %s -o output url", os.Args[0])
	}
	page := args[0]

	var stylesheet string
	if *s.Stylesheet != "" {
		stylesheet = *s.Stylesheet
	} else {
		stylesheet = filepath.Join(filepath.Dir(os.Args[0]), "lib/latex/stylesheet.yml")
	}
	converter, err := latex.New(stylesheet)
	if err != nil {
		return err
	}

	log.Printf("Parsing %s\n", page)

	resp, err := http.Get(page)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ctx := util.NewBuffers(context.Background())

	if *s.Output == "-" {
		ctx = util.WithContext(os.Stdout, ctx)
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
		ctx = util.WithContext(f, ctx)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	return converter.Handler().Do(doc, ctx)
}
