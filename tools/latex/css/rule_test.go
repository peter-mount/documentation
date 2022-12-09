package css

import (
	"fmt"
	"os"
	"testing"
)

func TestParseRule(t *testing.T) {

	rules := []string{
		"td:not(:last-child)",
		"td:not(:last-child),th:not(:last-child)",
		"td:not(:last-child), th:not(:last-child)",
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
