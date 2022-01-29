package chip

import (
  "fmt"
  "github.com/peter-mount/go-kernel/util/strings"
)

type Category map[string]map[string]*Definition

func NewCategory() *Category {
  m := make(Category)
  return &m
}

// Put adds a Definition to the Category.
// It returns true if the definition was added, false if an entry already exists.
func (c *Category) Put(d *Definition) bool {
  if d.Category == "" {
    d.Category = "Miscellaneous"
  }

  m, exists := (*c)[d.Category]
  if !exists {
    m = make(map[string]*Definition)
    (*c)[d.Category] = m
  }

  if _, exists = m[d.Name]; exists {
    return false
  }

  m[d.Name] = d
  return true
}

func (c *Category) Get(cat, name string) *Definition {
  if m, exists := (*c)[cat]; exists {
    return m[name]
  }
  return nil
}

// Categories returns a sorted slice of category names
func (c *Category) Categories() strings.StringSlice {
  var a strings.StringSlice

  for k, _ := range *c {
    a = append(a, k)
  }

  return a.Sort()
}

func (c *Category) DefinitionNames(cat string) strings.StringSlice {
  var a strings.StringSlice

  if m, exists := (*c)[cat]; exists {
    for k, _ := range m {
      a = append(a, k)
    }
  }

  return a.Sort()
}

func (c *Category) ForEachCategory(f func(string) error) error {
  return c.Categories().ForEach(f)
}

func (c *Category) ForEach(cat string, f func(string, *Definition) error) error {
  m, exists := (*c)[cat]
  if !exists {
    return fmt.Errorf("category %q unknown", cat)
  }

  return c.DefinitionNames(cat).
      ForEach(func(k string) error {
        return f(k, m[k])
      })
}
