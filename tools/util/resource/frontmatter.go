package resource

import (
  "github.com/peter-mount/documentation/tools/util"
  strings2 "github.com/peter-mount/go-kernel/util/strings"
  "path"
  "sort"
  "strings"
)

// Flatten returns the Resource & all children as a single slice.
// Children will have their names prefixed with their parents so the names
// appear as if they are in a subdirectory.
func (r *resource) Flatten() []Resource {
  var a []Resource
  a = r.resources(a, "")
  sort.SliceStable(a, func(i, j int) bool {
    return strings.ToLower(a[i].Name()) < strings.ToLower(a[j].Name())
  })
  return a
}

func (r *resource) resources(a []Resource, prefix string) []Resource {
  if r.Size() > 0 {
    a = append(a, Wrap(prefix, r))
  }

  _ = r.ForEach(func(c Resource) error {
    if cr, ok := c.(*resource); ok {
      np := r.Name()
      if prefix != "" {
        np = path.Join(prefix, np)
        if r.Size() > 0 {
          a = append(a, Wrap(prefix, r))
        }
      } else {
        a = cr.resources(a, path.Join(prefix, np))
      }
    }
    return nil
  })

  return a
}

// Wrap returns a resource with a directory prefix attached to it's name
func Wrap(prefix string, r Resource) Resource {
  return &resource{
    name: path.Join(prefix, r.Name()),
    url:  r.Url(),
    size: r.Size(),
  }
}

func (r resource) FileBuilder() util.FileBuilder {
  return func(slice strings2.StringSlice) (strings2.StringSlice, error) {
    res := r.Flatten()
    if len(res) > 0 {
      slice = append(slice, "resources:")
      for _, e := range res {
        slice = append(slice, "  - name: \""+e.Name()+"\"")
        slice = append(slice, "    url: \""+e.Url()+"\"")
        slice = append(slice, "    size: \""+util.Unit(e.Size())+"\"")
      }
    }
    return slice, nil
  }
}
