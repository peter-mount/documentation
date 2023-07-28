package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

const (
	StyleChapter              = `\chapter`
	StyleSection              = `\section`
	StyleSectionPaged         = `\sectionpaged`
	StyleSectionPagedLeft     = `\sectionpagedleft`
	StyleSectionPagedRight    = `\sectionpagedright`
	StyleSubSection           = `\subsection`
	StyleSubSectionPaged      = `\subsectionpaged`
	StyleSubSectionPagedLeft  = `\subsectionpagedleft`
	StyleSubSectionPagedRight = `\subsectionpagedright`
	StyleSubSubSection        = `\subsubsection`
	StyleSubSubSectionPaged   = `\subsubsectionpaged`
	StyleSubSubSectionLeft    = `\subsubsectionleft`
	StyleSubSubSectionRight   = `\subsubsectionright`
	StyleParagraph            = `\paragraph`
)

// heading handles headings h1 h2 etc. based on class type
func (c *Converter) headingStart(n *html.Node, ctx context.Context) error {
	var err error

	hType := n.Data

	// Get the style, note ordering is important here
	// as we want class to override the element name
	var style string
	switch {
	// title is for a non-numbered heading?
	case parser.HasClass(n, "chapter"),
		parser.HasClass(n, "title"):
		style = StyleChapter

	case parser.HasClass(n, "section"):
		style = StyleSection

	case parser.HasClass(n, "sectionpaged"):
		style = StyleSectionPaged

	case parser.HasClass(n, "sectionpagedleft"):
		style = StyleSectionPagedLeft

	case parser.HasClass(n, "sectionpagedright"):
		style = StyleSectionPagedRight

	case parser.HasClass(n, "subsection"):
		style = StyleSubSection

	case parser.HasClass(n, "subsectionpaged"):
		style = StyleSubSectionPaged

	case parser.HasClass(n, "subsectionpagedleft"):
		style = StyleSubSectionPagedLeft

	case parser.HasClass(n, "subsectionpagedright"):
		style = StyleSubSectionPagedRight

	case parser.HasClasses(n, "subsubsection", "subsubsectionbreak"):
		style = StyleSubSubSection

	case parser.HasClass(n, "subsubsectionpaged"):
		style = StyleSubSubSectionPaged

	case parser.HasClass(n, "subsubsectionleft"):
		style = StyleSubSubSectionLeft

	case parser.HasClass(n, "subsubsectionright"):
		style = StyleSubSubSectionRight

	case parser.HasClasses(n, "paragraph", "paragraphbreak"):
		style = StyleParagraph

	case hType == "h1":
		style = StyleSection

	case hType == "h2":
		style = StyleSubSection

	case hType == "h3":
		style = StyleSubSubSection

	case hType == "h4", hType == "h5":
		style = StyleParagraph

		// Default numbered heading
	default:
		style = StyleSection
	}
	err = util.WriteString(ctx, style)

	if err == nil {
		err = util.Write(ctx, '{')
	}

	return err
}

func (c *Converter) headingEnd(n *html.Node, ctx context.Context) error {
	if parser.HasClasses(n, "subsubsectionbreak", "paragraphbreak") {
		return util.WriteStringLn(ctx, `}\\~\\[-\baselineskip]`)
	}
	return util.Write(ctx, '}', '\n')
}
