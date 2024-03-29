package css

import (
	"fmt"
	"os"
	"testing"
)

var rules = []string{
	// Basic rule
	"td:not(:last-child)",
	// Rule's with 2 targets, one with whitespace
	"td:not(:last-child),th:not(:last-child)",
	"td:not(:last-child), th:not(:last-child)",
	// Rule which handles matching 2 elements in sequence
	// e.g. and(tbody td)
	"tbody td:first-child",
	// Similarly but with 3
	// e.g. and(and(tbody td) strong)
	"tbody td:first-child strong",
	// element.class selector
	"table.memory",
	// .class selector
	".memory",
}

func TestParseRule(t *testing.T) {

	for i, rs := range rules {
		t.Run(fmt.Sprintf("rule%d", i), func(t *testing.T) {
			r, err := ParseRule(rs)
			if err != nil {
				t.Error(err)
			} else {
				r.Write(os.Stdout)
			}
		})
	}
}
