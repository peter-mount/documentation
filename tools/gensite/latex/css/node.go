package css

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/gensite/latex/util"
	"github.com/tdewolff/parse/css"
	"io"
	"strings"
)

type NodeAction func(*Context, *Node) (*util.Value, error)

type Node struct {
	Parent    *Node         // Parent node
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
	NodeAnd      = "and"     // Special function node which ands two results, e.g. "tr td" is 3 nodes: and(tr td)
	NodeClass    = "class"   // Class selector
)

func (n *Node) Do(ctx *Context) (*util.Value, error) {
	if n == nil || n.Action == nil {
		return nil, nil
	}

	return n.Action(ctx, n)
}

// Root returns the root node in the tree
func (n *Node) Root() *Node {
	if n == nil || n.Parent == nil {
		return n
	}
	return n.Parent.Root()
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

	b.Parent = n
}

func (n *Node) parse(l *css.Lexer) error {
	for {
		tt, text := l.Next()
		//fmt.Printf("np: %v %q\n", tt, string(text))
		switch tt {
		case css.ErrorToken:
			// error or EOF set in l.Err()
			return l.Err()

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

		case css.WhitespaceToken:
			// Whitespace is an and operation between two rules
			root := &Node{Type: NodeAnd, TokenType: css.FunctionToken}
			root.add(n.Root())

			// Start new run on RHS
			rhs, err := parseIdent(l)
			if err != nil {
				return err
			}
			root.add(rhs)

			// Finish off parsing the new RHS
			return rhs.parse(l)

			// Ignore colons, we use them as delimiters here
		case css.ColonToken:

		case css.DelimToken:
			// Class delimiter
			root, err := n.parseClass(l)
			if err != nil {
				return err
			}
			root.Action = delim

			// Finish off by carrying on but this will add to rhs with next tokens
			return root.parse(l)

		default:
			fmt.Printf("unk: %v %q\n", tt, string(text))
		}
	}
}

func (n *Node) parseClass(l *css.Lexer) (*Node, error) {
	root := &Node{Type: NodeClass, TokenType: css.DelimToken}
	root.add(n)

	// Start new run on RHS
	rhs, err := parseIdent(l)
	if err != nil {
		return nil, err
	}

	// We don't add it but copy its text as that's the class name
	root.Text = rhs.Text
	return root, nil
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
