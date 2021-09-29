package util

import (
  "go/scanner"
  "go/token"
)

type Token struct {
  Next   *Token      // Next token in chain
  Child  *Token      // First child token
  Parent *Token      // Parent token
  Token  token.Token // Token
  Lit    string      // Literal
}

func (t *Token) next(nt *Token) *Token {
  if t != nil {
    t.Next = nt
  }
  return nt
}

func (t *Token) child(nt *Token) *Token {
  if t != nil {
    t.Child = nt
    nt.Parent = t
  }
  return nt
}

func (t *Token) String() string {
  if t == nil {
    return "nil"
  }
  if t.Lit != "" {
    return t.Lit
  }
  return t.Token.String()
}

func tokenize(l int, s *scanner.Scanner) *Token {
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
        _ = t.child(tokenize(l+1, s))
      } else {
        // Already have a child so add LPAREN for ()
        nt = t.next(&Token{Token: token.LPAREN}).
          child(tokenize(l+1, s))
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

    //fmt.Printf("tok %s lit %q child=%x\n", t.Token, t.Lit, t.Child)
  }
}

func Tokenize(src string) *Token {
  var s scanner.Scanner
  fSet := token.NewFileSet()
  file := fSet.AddFile("", fSet.Base(), len(src))
  s.Init(file, []byte(src), nil, scanner.ScanComments)
  return tokenize(0, &s)
}
