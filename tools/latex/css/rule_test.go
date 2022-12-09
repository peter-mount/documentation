package css

import (
	"fmt"
	"os"
	"testing"
)

func TestParseRule(t *testing.T) {

	rules := []string{
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
	}
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
