package autodoc

import (
  "sort"
  "strings"
)

type ApiEntry struct {
  call     int            // Call 0..255
  Name     string         `yaml:"name"`
  Addr     string         `yaml:"addr"`
  Indirect string         `yaml:"indirect,omitempty"`
  Title    string         `yaml:"title"`
  Entry    FunctionParams `yaml:"entry"`
  Exit     FunctionParams `yaml:"exit"`
  params   interface{}
}

type ApiEntryHandler func(*ApiEntry) error

func (a *Api) SortByName() {
  sort.SliceStable(a.api, func(i, j int) bool {
    return strings.ToLower(a.api[i].Name) < strings.ToLower(a.api[j].Name)
  })
}

func (a *Api) SortByAddr() {
  sort.SliceStable(a.api, func(i, j int) bool {
    return a.api[i].call < a.api[j].call
  })
}
