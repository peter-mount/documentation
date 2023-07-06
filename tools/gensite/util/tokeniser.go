package util

import (
	"go/scanner"
	"go/token"
)

// Token represents a parsed token as a linked list.
// If a () is found then it represents a child to the previous token, unless the previous token was also a ()
// in which case it's taken to be a scanner.LPAREN token in its own right with the contents being its child.
type Token struct {
	Next   *Token      // Next token in chain
	Child  *Token      // First child token
	Parent *Token      // Parent token
	Token  token.Token // Token
	Lit    string      // Literal
}

// next Adds a token to the current one returning the new token
func (t *Token) next(nt *Token) *Token {
	if t != nil {
		t.Next = nt
	}
	return nt
}

// child adds a token as a child, returning the new token
func (t *Token) child(nt *Token) *Token {
	if t != nil && nt != nil {
		t.Child = nt
		nt.Parent = t
	}
	return nt
}

// String converts this token to a string
func (t *Token) String() string {
	if t == nil {
		return "nil"
	}
	if t.Lit != "" {
		return t.Lit
	}
	return t.Token.String()
}

// Tokenize a string returning a Token tree.
//
// Most tokens form a linked list, however if a token.LPAREN is found then the tokens after it until a matching
// token.RPAREN are placed as children of the current token. This allows for simple nesting of functions.
//
// If a token.LPAREN follows a token.RPAREN, then a special token is appended of type token.LPAREN to represent that it
// contains child tokens with no parent.
func Tokenize(src string) *Token {
	var s scanner.Scanner
	fSet := token.NewFileSet()
	file := fSet.AddFile("", fSet.Base(), len(src))
	s.Init(file, []byte(src), nil, scanner.ScanComments)
	return tokenize(&s)
}

// tokenize accepts tokens from the scanner.
// If token.EOF then it terminates.
// If token.LPAREN it then recurse using the current token as the parent, implementing nesting.
// If token.RPAREN it returns to the caller, usually the parent token.
// otherwise, it adds the token as the next sibling to the current token.
func tokenize(s *scanner.Scanner) *Token {
	var t *Token
	var rt *Token
	for {
		var nt *Token
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			return rt
		case token.LPAREN:
			if t != nil && t.Child == nil {
				// Add child but ignore the returned token as we want to stay
				// on this one
				_ = t.child(tokenize(s))
			} else {
				// Already have a child so add LPAREN for ()
				nt = t.next(&Token{Token: token.LPAREN}).
					child(tokenize(s))
			}
		case token.RPAREN:
			return rt
		default:
			nt = t.next(&Token{Token: tok, Lit: lit})
		}

		if nt != nil {
			t = nt
			if rt == nil {
				rt = nt
			}
		}
	}
}
