package pdf

import (
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/peter-mount/documentation/tools/gendoc"
	"github.com/peter-mount/documentation/tools/gendoc/hugo"
	"github.com/peter-mount/documentation/tools/gendoc/util"
	"github.com/peter-mount/documentation/tools/gendoc/web"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"os"
	"strings"
)

// PDF tool that handles the generation of PDF documentation of a "book"
type PDF struct {
	config    *Config         `kernel:"config,pdf"` // Config
	bookShelf *hugo.BookShelf `kernel:"inject"`     // Bookshelf
	chromium  *web.Chromium   `kernel:"inject"`     // Chromium browser
	worker    task.Queue      `kernel:"worker"`     // Worker queue
	webserver *web.Webserver  `kernel:"inject"`     // Webserver
	_         *hugo.Hugo      `kernel:"inject"`     // access them directly
}

type Config struct {
	PrintBackground     bool    `yaml:"printBackground"`     // Print background graphics
	Margin              Margin  `yaml:"margin"`              // Page margins
	Width               float64 `yaml:"width"`               // Page width in inches
	Height              float64 `yaml:"height"`              // Page height in inches
	Landscape           bool    `yaml:"landscape"`           // Landscape or Portrait
	Header              string  `yaml:"header"`              // Header in html
	Footer              string  `yaml:"footer"`              // Footer in html
	DisableHeaderFooter bool    `yaml:"disableHeaderFooter"` // Disable header & footer generation
}

// Margin in inches. Default is 1cm ~0.4 inches
type Margin struct {
	Top    float64 `yaml:"top"`    // Top margin in inches
	Left   float64 `yaml:"left"`   // Left margin in inches
	Bottom float64 `yaml:"bottom"` // Bottom margin in inches
	Right  float64 `yaml:"right"`  // Right margin in inches
}

func (p *PDF) Start() error {
	args := flag.Args()
	if len(args) != 2 {
		return fmt.Errorf("syntax: %s bookId dest.pdf", os.Args[0])
	}
	bookId, destFile := args[0], args[1]

	return p.bookShelf.
		Books().
		Filter(func(book *hugo.Book) bool {
			return book.ID == bookId
		}).
		ForEach(func(book *hugo.Book) error {
			p.worker.AddPriorityTask(tools.PriorityPDF, func(ctx context.Context) error {
				log.Println("Generating PDF for", book.ID)
				return p.chromium.Run(p.printToPDF(book, destFile))
			})
			return nil
		})
}

// print a specific pdf page.
func (p *PDF) printToPDF(book *hugo.Book, destFile string) chromedp.Tasks {
	url := p.webserver.WebPath("%s/_print/", strings.ToLower(book.WebPath()))

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(p.config.PrintBackground).
				WithMarginTop(p.config.Margin.Top).
				WithMarginBottom(p.config.Margin.Bottom).
				WithMarginLeft(p.config.Margin.Left).
				WithMarginRight(p.config.Margin.Right).
				WithLandscape(p.config.Landscape).
				WithPaperWidth(p.config.Width).
				WithPaperHeight(p.config.Height).
				WithDisplayHeaderFooter(!p.config.DisableHeaderFooter).
				WithHeaderTemplate(book.Expand(p.config.Header)).
				WithFooterTemplate(book.Expand(p.config.Footer)).
				Do(ctx)

			if err != nil {
				return err
			}

			return util.ByteFileHandler(buf).
				Write(destFile, book.Modified())
			//Write("public/static/book/"+book.ID+".pdf", book.Modified())
		}),
	}
}
