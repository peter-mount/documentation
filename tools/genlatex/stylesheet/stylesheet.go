package stylesheet

// Stylesheet is read from yaml and defines how certain html elements
// are handled based on the declared classes
type Stylesheet struct {
	Table        map[string]*Table `yaml:"table"` // Table definitions
	defaultTable *Table            `yaml:"-"`     // default table
}

// Init initialises the Stylesheet with defaults
func (s *Stylesheet) Init() {

	// Initialise tables so they have defaults
	s.defaultTable = &Table{}
	s.defaultTable.init()
	for _, t := range s.Table {
		t.init()
	}
}

func defString(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
