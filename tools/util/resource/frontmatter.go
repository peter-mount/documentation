package resource

import (
  "github.com/peter-mount/documentation/tools/util"
  "path"
  "sort"
  "strings"
)

// Flatten returns the Resource & all children as a single slice.
// Children will have their names prefixed with their parents so the names
// appear as if they are in a subdirectory.
func (r *resource) Flatten() []Resource {
  var a []Resource
  a = r.resources(a, r.Name())
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
      a = cr.resources(a, path.Join(prefix, r.Name()))
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
  return func(slice util.StringSlice) (util.StringSlice, error) {
    slice = append(slice, "resources:")
    for _, e := range r.Flatten() {
      slice = append(slice, "  - name: \""+e.Name()+"\"")
      slice = append(slice, "    url: \""+e.Url()+"\"")
      slice = append(slice, "    size: \""+util.Unit(e.Size())+"\"")
    }
    return slice, nil
  }
}
