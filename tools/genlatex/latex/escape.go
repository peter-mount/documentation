package latex

import "strings"

var (
	escapeFrom []string
	escapeTo   []string
)

func init() {
	initEscape := func(from, to string) {
		escapeFrom = append(escapeFrom, from)
		escapeTo = append(escapeTo, to)
	}

	// Note {} after each \command so that LaTeX knows the end of each
	// one and it doesn't then either get confused or add unwanted spaces.
	// e.g. "shift-lock" becomes "shift\textendashlock" and fails but
	// works with "shift\textendash{}lock"

	// TODO look at how to handle the existing latex equations on the site
	// as this will break those

	// These must be first otherwise it will break all that follow
	initEscape(`{`, `\{`)
	initEscape(`}`, `\}`)
	initEscape(`\`, `\textbackslash{}`)

	// ~ is either \~{} or 	extasciitilde{}
	// Must be before &nbsp;
	initEscape(`~`, `\textasciitilde{}`)

	// html characters
	initEscape(`&nbsp;`, `~`)

	initEscape(`&asymp;`, `\AA`)
	initEscape(`≈`, `\AA`)

	initEscape(`&aring;`, `$\approx$`)
	initEscape(`Å`, `$\approx$`)

	initEscape(`&ne;`, `$\neq$`)
	initEscape(`≠`, `$\neq$`)

	initEscape(`&times;`, `$\times$`)
	initEscape(`×`, `$\times$`)

	initEscape(`&middot;`, `$\cdot$`)
	initEscape(`·`, `$\cdot$`)

	initEscape(`&divide;`, `$\div$`)
	initEscape(`÷`, `$\div$`)

	initEscape(`&lt;`, `\textless{}`)
	initEscape(`<`, `\textless{}`)

	initEscape(`&le;`, `$\leq$`)
	initEscape(`≤`, `$\leq$`)

	initEscape(`&gt;`, `\textgreater{}`)
	initEscape(`>`, `\textgreater{}`)

	initEscape(`&ge;`, `$\geq$`)
	initEscape(`≥`, `$\geq$`)

	initEscape(`&radic;`, `$\sqrt{}$`)
	initEscape(`√`, `$\sqrt{}$`)

	initEscape(`&sup2;`, `$^2$`)
	initEscape(`²`, `$^2$`)

	initEscape(`&sup3;`, `$^3$`)
	initEscape(`³`, `$^3$`)

	initEscape(`&deg;`, `$^\circ$`)
	initEscape(`°`, `$^\circ$`)

	initEscape(`&micro;`, `$\mu$`)
	initEscape(`µ`, `$\mu$`)

	// general symbols
	initEscape(`%`, `\%`)
	initEscape(`$`, `\$`)
	initEscape(`_`, `\_`)
	initEscape(`#`, `\#`)
	initEscape(`&`, `\&`)
	initEscape(`|`, `\textbar{}`)
	initEscape(`‡`, `\ddag{}`)
	initEscape(`-`, `\textendash{}`)
	initEscape(`£`, `\pounds{}`)
	initEscape(`†`, `\dag{}`)
}

// EscapeText takes a string and converts it to LaTeX
func EscapeText(s string) string {
	for i, v := range escapeFrom {
		s = strings.ReplaceAll(s, v, escapeTo[i])
	}
	return s
}
