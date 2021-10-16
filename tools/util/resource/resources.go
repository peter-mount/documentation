package resource

// Resource represents the Resources table on the top right side of each page
// It can represent either a single downloadable resource or a directory of
// resources.
type Resource interface {
  // Name of this Resource
  Name() string
  // Url of this Resource
  Url() string
  // Size of this Resource in bytes
  Size() int
  // Add a Resource to to this instance
  Add(name, url string, size int) Resource
  // AddChild adds a sub-Resource to this one, e.g. a directory
  AddChild(Resource) Resource
  // Flatten returns the Resource & all children as a single slice.
  Flatten() []Resource
  // ForEach invokes a Handler for each Resource including children
  ForEach(Handler) error
}

type Handler func(Resource) error

type resource struct {
  name     string     // Name of the resource
  url      string     // Path to the resource
  size     int        // Size of the resource
  children []Resource // Child resources
}

// Add a specific resource to this instance
func (r *resource) Add(name, url string, size int) Resource {
  return r.AddChild(&resource{name: name, url: url, size: size})
}

// AddChild adds a child Resource, e.g. resources for a sub-page we want to include here
func (r *resource) AddChild(child Resource) Resource {
  r.children = append(r.children, child)
  return r
}

func (r *resource) ForEach(f Handler) error {
  for _, c := range r.children {
    err := f(c)
    if err != nil {
      return err
    }
  }
  return nil
}

// Name of this specific resource
func (r *resource) Name() string {
  return r.name
}

// Url to this resource online
func (r *resource) Url() string {
  return r.url
}

// Size of this resource in bytes
func (r *resource) Size() int {
  return r.size
}
