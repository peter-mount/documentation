package autodoc


// Output is used for generating the index pages front matter
type Output struct {
  Nometa bool          `yaml:"nometa"`
  Api    []interface{} `yaml:"api,omitempty"`
}

