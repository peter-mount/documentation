package chip

import (
	"github.com/peter-mount/documentation/tools/gensite/util"
	"github.com/peter-mount/documentation/tools/gensite/util/html"
	"go/token"
)

type PinFormatter struct {
	e              *html.Element
	fontSize       int
	textDecoration string
	dy             int
}

func Parse(s string, fontSize int, rightAlign bool, e *html.Element) *html.Element {
	f := &PinFormatter{
		e:        e.SvgText(),
		fontSize: fontSize,
	}

	if rightAlign {
		f.e = f.e.Attr("text-anchor", "end")
	}

	f.parse(util.Tokenize(s))
	return e
}

func (f *PinFormatter) subscript(t *util.Token) {
	curE := f.e
	curFontSize := f.fontSize
	curTD := f.textDecoration
	defer func() {
		f.fontSize = curFontSize
		f.e = curE
		f.textDecoration = curTD
	}()
	f.textDecoration = ""
	if f.fontSize > 6 {
		f.fontSize = f.fontSize - 1
	}

	// Next text() will offset by this
	f.dy = 5

	f.e = f.e.TSpan().
		AttrInt("font-size", f.fontSize).
		AttrInt("dx", -1)
	f.parse(t)

	// If f.dy was reset then set it to inverse to adjust the next text() back inline
	if f.dy == 0 {
		f.dy = -5
	}
}

func (f *PinFormatter) text(s string) {
	f.e = f.e.TSpan()
	if f.textDecoration != "" {
		f.e.Class(f.textDecoration)
	}
	if f.dy != 0 {
		f.e.AttrInt("dy", f.dy)
		f.dy = 0
	}
	f.e = f.e.Text(s).End()
}

func (f *PinFormatter) reference(l string, t *util.Token) {
	curE := f.e
	defer func() {
		f.e = curE
	}()
	f.text(l)
	f.subscript(t.Child)
}

func (f *PinFormatter) parse(t *util.Token) {
	curE := f.e
	curTD := f.textDecoration
	defer func() {
		f.e = curE
		f.textDecoration = curTD
	}()

	for t != nil {
		if t.Child != nil {
			switch t.Token {
			case token.SEMICOLON:
				// Ignore
			case token.LPAREN:
				f.text(" (")
				f.parse(t.Child)
				f.text(")")

			case token.IDENT:
				switch t.Lit {
				// Text with a line above it
				case "NOT":
					f.textDecoration = "not"
					f.parse(t.Child)
				case "PHI":
					f.reference("∅", t)
				// Default the ident as text then expand the child content
				default:
					f.reference(t.String(), t)
				}
			// Default the ident as text then expand the child content
			default:
				f.reference(t.String(), t)
			}
		} else if t.Token != token.SEMICOLON {
			// For all content other than token.SEMICOLON just append it as text
			f.text(t.String())
		}
		t = t.Next
	}
}
