package hugo

import (
	"bufio"
	"context"
	"github.com/peter-mount/documentation/tools/gensite/util"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type FrontMatter struct {
	Type       string                 `yaml:"type"`       // Template for this page
	Title      string                 `yaml:"title"`      // Title of page
	LinkTitle  string                 `yaml:"linkTitle"`  // Link title to this page
	Weight     int                    `yaml:"weight"`     // Weight, 0=none. Use -1 if you specifically need 0
	Categories []string               `yaml:"categories"` // Page categories
	Tags       []string               `yaml:"tags"`       // Page tags
	Book       *Book                  `yaml:"book"`       // Book definition
	Other      map[string]interface{} `yaml:",inline"`    // All other data in raw format
}

type FrontMatterAction func(context.Context, *FrontMatter) error

func FrontMatterActionOf(actions ...FrontMatterAction) FrontMatterAction {
	switch len(actions) {
	case 0:
		return nil
	case 1:
		return actions[0]
	default:
		a := actions[0]
		for _, b := range actions[1:] {
			a = a.Then(b)
		}
		return a
	}
}

func (a FrontMatterAction) Then(b FrontMatterAction) FrontMatterAction {
	if a == nil {
		return b
	}
	return func(ctx context.Context, fm *FrontMatter) error {
		if err := a.Do(ctx, fm); err != nil {
			return err
		}
		return b.Do(ctx, fm)
	}
}

func (a FrontMatterAction) OtherExists(key string, f FrontMatterAction) FrontMatterAction {
	return a.Then(func(ctx context.Context, fm *FrontMatter) error {
		if val, exists := fm.Other[key]; exists {
			return f.Do(context.WithValue(ctx, "other", val), fm)
		}
		return nil
	})
}

func (a FrontMatterAction) WithNotes(globalNotes *util.Notes) FrontMatterAction {
	return func(ctx context.Context, fm *FrontMatter) error {
		ctx = context.WithValue(ctx, "globalNotes", globalNotes)

		notes := util.NewNotes()
		notes.DecodePageNotes(fm.Other["notes"])
		ctx = context.WithValue(ctx, "notes", notes)
		defer globalNotes.Merge(notes)

		return a.Do(ctx, fm)
	}
}

func (a FrontMatterAction) Context(name string, value interface{}) FrontMatterAction {
	if a == nil {
		return nil
	}

	return func(ctx context.Context, matter *FrontMatter) error {
		return a(context.WithValue(ctx, name, value), matter)
	}
}

func (a FrontMatterAction) Do(ctx context.Context, fm *FrontMatter) error {
	if a == nil {
		return nil
	}
	return a(ctx, fm)
}

func (a FrontMatterAction) Walk(ctx context.Context) walk.PathWalker {
	return func(path string, fileInfo os.FileInfo) error {
		fm := &FrontMatter{}
		if err := fm.LoadFrontMatter(path); err != nil {
			return err
		}

		ctx = context.WithValue(ctx, "fileInfo", fileInfo)
		return a.Do(ctx, fm)
	}
}

func FileInfo(ctx context.Context) os.FileInfo {
	fi, ok := ctx.Value("fileInfo").(os.FileInfo)
	if ok {
		return fi
	}
	return nil
}

func (fm *FrontMatter) LoadFrontMatter(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return fm.ReadFrontMatter(f)
}

func (fm *FrontMatter) ReadFrontMatter(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		if l == "---" {
			return fm.readFrontMatter(scanner)
		}
	}

	// No frontmatter
	return nil
}

func (fm *FrontMatter) readFrontMatter(s *bufio.Scanner) error {
	var a []string
	for s.Scan() {
		l := s.Text()
		if l == "---" {
			b := []byte(strings.Join(a, "\n"))
			return yaml.Unmarshal(b, fm)
		}

		a = append(a, l)
	}

	// no --- terminator so ignore the page
	return nil
}
