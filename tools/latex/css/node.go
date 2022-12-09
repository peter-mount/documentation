package css

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/latex/util"
	"github.com/tdewolff/parse/css"
	"io"
	"strings"
)

type NodeAction func(context.Context, *Node) (*util.Value, error)

type Node struct {
	Left      *Node         // Left hand side
	Right     *Node         // Right hand side
	Text      string        // Text for the comparison
	Type      NodeType      // Node type
	TokenType css.TokenType // TokenType
	Action    NodeAction    // Action to perform
}

type NodeType string

const (
	NodeElement  = "element" // Select against the named element type
	NodeFunction = "func"    // Function call
)

func (n *Node) Do(ctx context.Context) (*util.Value, error) {
	if n == nil || n.Action == nil {
		return nil, nil
	}

	return n.Action(ctx, n)
}

func (n *Node) Write(w io.Writer, i int) {
	s := strings.Repeat(" ", i) + "|"
	_, _ = fmt.Fprintf(w, s+"tkn: %s %s %q\n", n.Type, n.TokenType, n.Text)
	if n.Left != nil {
		_, _ = fmt.Fprintf(w, s+"lhs:\n")
		n.Left.Write(w, i+2)
	}
	if n.Right != nil {
		_, _ = fmt.Fprintf(w, s+"rhs:\n")
		n.Right.Write(w, i+2)
	}
}

func (n *Node) add(b *Node) {
	if n == b {
		panic("Cannot add node to itself")
	}
	switch {
	case n.Left == nil:
		n.Left = b
	case n.Right == nil:
		n.Right = b
	default:
		panic("Node full")
	}
}

func (n *Node) parse(l *css.Lexer) error {
	for {
		tt, text := l.Next()
		//fmt.Printf("np: %v %q\n", tt, string(text))
		switch tt {
		case css.ErrorToken:
			// error or EOF set in l.Err()
			return l.Err()
		case css.ColonToken, css.WhitespaceToken:
		// Function call
		case css.FunctionToken:
			nn, err := n.parseFunction(l, string(text))
			if err != nil {
				return err
			}
			n.add(nn)
			// Valid for functions only
		case css.RightParenthesisToken:
			return nil

		case css.CommaToken:
			// End the node processing run
			return nil

		case css.IdentToken:
			nn, err := n.parseIdent(l, string(text))
			if err != nil {
				return err
			}
			n.add(nn)

		default:
			//fmt.Printf("unk: %v %q\n", tt, string(text))
		}
	}
}

func (n *Node) parseFunction(l *css.Lexer, name string) (*Node, error) {
	return n.parseFuncIdent(l, name, true)
}

func (n *Node) parseIdent(l *css.Lexer, name string) (*Node, error) {
	return n.parseFuncIdent(l, name, false)
}

func (n *Node) parseFuncIdent(l *css.Lexer, name string, recurse bool) (*Node, error) {
	name = strings.TrimSuffix(name, "(")

	action, ok := funcMap[name]
	if !ok {
		fmt.Printf("wrn: Unknown func %q\n", name)
	}

	nn := &Node{Text: name, TokenType: css.FunctionToken, Type: NodeFunction, Action: action}

	if recurse {
		err := nn.parse(l)
		if err != nil {
			return nil, err
		}
	}

	return nn, nil
}
