package bbc

import (
	"flag"
	"github.com/peter-mount/documentation/tools/util"
	"github.com/peter-mount/go-kernel"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type BBC struct {
	baseDir *string
	api     []*Api
	osbyte  []*Osbyte
	osword  []*Osword
}

type Output struct {
	Nometa bool          `yaml:"nometa"`
	Api    []interface{} `yaml:"api,omitempty"`
	Osbyte []interface{} `yaml:"osbyte,omitempty"`
	Osword []interface{} `yaml:"osword,omitempty"`
}

func (b *BBC) Name() string {
	return "BBC"
}

func (b *BBC) Init(_ *kernel.Kernel) error {
	b.baseDir = flag.String("bbc", "", "Src dir for bbc content")

	return nil
}

func (b *BBC) Run() error {
	// Do nothing
	if *b.baseDir == "" {
		return nil
	}

	err := b.extract()
	if err != nil {
		return err
	}

	err = b.writeAPIIndex()
	if err != nil {
		return err
	}

	err = b.writeAPINameIndex()
	if err != nil {
		return err
	}

	err = b.writeOsbyteIndex()
	if err != nil {
		return err
	}

	err = b.writeOswordIndex()
	if err != nil {
		return err
	}

	return nil
}

func (b *BBC) extract() error {
	b.osbyte = nil

	log.Println("Scanning BBC API")
	err := filepath.Walk(*b.baseDir, b.extractMeta)
	if err != nil {
		return err
	}

	return nil
}

func (b *BBC) extractMeta(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") && !strings.Contains(path, "reference") {
		//log.Println("Reading", path)
		m, err := util.LoadFrontMatter("", path)
		if err != nil {
			return err
		}

		if api, exists := m["api"]; exists {
			b.extractApi(api)
		}

		if osbyte, exists := m["osbyte"]; exists {
			b.extractOsbyte(osbyte)
		}

		if osword, exists := m["osword"]; exists {
			b.extractOsword(osword)
		}
	}

	return nil
}
