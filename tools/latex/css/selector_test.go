package css

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/latex/util"
  "github.com/tdewolff/parse/css"
  "golang.org/x/net/html"
  "gopkg.in/yaml.v2"
  "io"
  "log"
  "os"
  "strings"
  "testing"
)

const (
  styleCss = "" +
      "styles:\n" +
      "  - rule: \"table.memory\"\n" +
      "    children:\n" +
      "      - rule: \"tr\"\n" +
      "        css:\n" +
      "          border-bottom: \"1px solid lightgrey\"\n" +
      "      - rule: \"th:not(:last-child), td:not(:last-child)\"\n" +
      "        css:\n" +
      "          border-right: \"1px solid lightgrey\"\n" +
      "          vertical-align: \"top\"\n" +
      "      - rule: \"th, td\"\n" +
      "        css:\n" +
      "          vertical-align: \"top\""
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

  err = Parse(styles.Styles)
  if err != nil {
    return nil, nil, err
  }

  Write(os.Stdout, styles.Styles)

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

func TestStyle_Parse(t *testing.T) {
  styles, _, err := testStyleSetup()
  if err != nil {
    t.Fatal(err)
  }

  for _, s := range styles.Styles[0].Children {
    log.Println(s.RuleSrc)
    for k, v := range s.Css {
      c := fmt.Sprintf("%s:%s", k, v)
      log.Println(c)

      l := css.NewLexer(strings.NewReader(c))
      for {
        tt, text := l.Next()
        log.Println(tt, text)
        switch tt {
        case css.ErrorToken:
          if l.Err() != io.EOF {
            fmt.Println("Error on line", l.Err())
          }
          return
        case css.IdentToken:
          fmt.Println("Identifier", string(text))
        case css.NumberToken:
          fmt.Println("Number", string(text))
          // ...
        }
      }
    }
  }

}
