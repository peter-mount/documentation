package hugo

import (
  "bufio"
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
