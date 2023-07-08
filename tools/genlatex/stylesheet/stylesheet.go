package stylesheet

// Stylesheet is read from yaml and defines how certain html elements
// are handled based on the declared classes
type Stylesheet struct {
	DocumentClass string            `yaml:"documentClass"`
	UsePackage    []string          `yaml:"usePackage"`
	Styles        map[string]Style  `yaml:"styles"` // Generic styles
	Table         map[string]*Table `yaml:"table"`  // Table definitions
	defaultTable  *Table            `yaml:"-"`      // default table
}

// Init initialises the Stylesheet with defaults
func (s *Stylesheet) Init() {
	s.DocumentClass = DefString(s.DocumentClass, "textbook")

	// Initialise tables so they have defaults
	s.defaultTable = &Table{}
	s.defaultTable.init()
	for _, t := range s.Table {
		t.init()
	}

	// For debugging the loaded stylesheet
	//b, _ := yaml.Marshal(s)
	//fmt.Println(string(b))
}

func DefString(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func DefStrings(a ...string) string {
	for _, s := range a {
		if s != "" {
			return s
		}
	}
	return ""
}
