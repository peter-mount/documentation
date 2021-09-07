package hugo

import (
  "bufio"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/walk"
  "gopkg.in/yaml.v2"
  "io"
  "os"
  "path"
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

type FrontMatterAction func(*FrontMatter) error

type FrontMatterOtherAction func(*FrontMatter, interface{}) error

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
  return func(fm *FrontMatter) error {
    if err := a.Do(fm); err != nil {
      return err
    }
    return b.Do(fm)
  }
}

func (a FrontMatterAction) OtherExists(key string, f FrontMatterOtherAction) FrontMatterAction {
  return a.Then(func(fm *FrontMatter) error {
    if val, exists := fm.Other[key]; exists {
      return f(fm, val)
    }
    return nil
  })
}

func (a FrontMatterAction) Do(fm *FrontMatter) error {
  if a == nil {
    return nil
  }
  return a(fm)
}

func (a FrontMatterAction) Walk() walk.PathWalker {
  return func(path string, _ os.FileInfo) error {
    fm := &FrontMatter{}
    if err := fm.LoadFrontMatter(path); err != nil {
      return err
    }

    return a.Do(fm)
  }
}

func (fm *FrontMatter) LoadFrontMatter(elem ...string) error {
  f, err := os.Open(path.Join(elem...))
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

func (a FrontMatterAction) WithNotes(global *util.Notes, f func(*util.Notes, *FrontMatter) error) FrontMatterAction {
  return a.Then(func(fm *FrontMatter) error {
    notes := util.NewNotes()
    notes.DecodePageNotes(fm.Other["notes"])

    if err := f(notes, fm); err != nil {
      return err
    }

    // Import these notes into the global pool
    global.Merge(notes)

    return nil
  })
}
