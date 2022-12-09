package css

import (
	"github.com/peter-mount/documentation/tools/latex/util"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"testing"
)

const (
	styleCss = "" +
		"styles:\n" +
		"  \"memory\":\n" +
		"    - rule: \"tr\"\n" +
		"      css:\n" +
		"        border-bottom: \"1px solid lightgrey\"\n" +
		"    - rule: \"th:not(:last-child), td:not(:last-child)\"\n" +
		"      css:\n" +
		"        border-right: \"1px solid lightgrey\"\n" +
		"        vertical-align: \"top\"\n" +
		"    - rule: \"th, td\"\n" +
		"      css:\n" +
		"        vertical-align: \"top\""
	testHtml = "<body><table class='memory'>" +
		"<thead><tr><th>head1</th><th>head2</th></tr></thead>" +
		"<tbody>" +
		"<tr><td>r1c1</td><td>r1c2</td></tr>" +
		"<tr><td>r2c1</td><td>r2c2</td></tr>" +
		"</tbody>" +
		"</table></body>"
)

func testStyleSetup() (*Styles, *html.Node, error) {
	var styles Styles
	err := yaml.Unmarshal([]byte(styleCss), &styles)
	if err != nil {
		return nil, nil, err
	}

	for _, s1 := range styles.Styles {
		for _, s2 := range s1 {
			r, err := ParseRule(s2.RuleSrc)
			if err != nil {
				return nil, nil, err
			}
			s2.Rule = r
		}
	}

	doc, err := html.Parse(strings.NewReader(testHtml))
	if err != nil {
		return nil, nil, err
	}

	return &styles, doc, nil
}

func TestStyle_Search(t *testing.T) {
	styles, doc, err := testStyleSetup()
	if err != nil {
		t.Fatal(err)
	}

	err = styles.Search(doc)
	if err != nil {
		t.Error(err)
	}

	err = util.DebugHtml(doc, os.Stdout)
	if err != nil {
		t.Error(err)
	}
}
