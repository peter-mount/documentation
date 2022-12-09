package css

import (
	"context"
	"fmt"
	"github.com/tdewolff/parse/css"
	"io"
	"strings"
)

type Rule struct {
	Nodes  []*Node    // Nodes in this rule
	Action RuleAction // Action to perform
}

type RuleAction func(context.Context, *Rule) error

func (r *Rule) Write(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Rule: %d nodes\n", len(r.Nodes))
	for i, n := range r.Nodes {
		_, _ = fmt.Fprintf(w, "Rule: %d/%d\n", i+1, len(r.Nodes))
		n.Write(w, 2)
	}
}

func ParseRule(s string) (*Rule, error) {
	l := css.NewLexer(strings.NewReader(s))
	r := &Rule{}
	for {
		n, err := parseIdent(l)

		if err == nil {
			err = n.parse(l)
		}

		if err == nil || err == io.EOF {
			// Add the actual root to the node list
			r.Nodes = append(r.Nodes, n.Root())
		}

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return r, nil
		}

	}
}

func parseIdent(l *css.Lexer) (*Node, error) {
	for {
		tt, text := l.Next()
		//fmt.Printf("rp: %v %q\n", tt, string(text))
		switch tt {
		case css.ErrorToken:
			// error or EOF set in l.Err()
			return nil, l.Err()
		case css.IdentToken:
			n := &Node{Text: string(text), TokenType: tt, Type: NodeElement}
			return n, nil
		default:
			//fmt.Printf("unk1: %v %q\n", tt, string(text))
		}
	}
}
