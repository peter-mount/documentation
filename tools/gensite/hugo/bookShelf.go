package hugo

import (
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"os"
	"path"
	"strings"
)

// BookShelf manages all Book's.
type BookShelf struct {
	books Books
}

func (bs *BookShelf) Start() error {
	log.Println("Searching for books")
	return walk.NewPathWalker().
		IsFile().
		PathNotContain("/reference/").
		PathHasSuffix("/_index.html").
		Then(bs.scanPage).
		Walk("content")
}

func (bs *BookShelf) scanPage(pathName string, _ os.FileInfo) error {
	fm := FrontMatter{}
	err := fm.LoadFrontMatter(pathName)
	if err != nil {
		return err
	}

	if fm.Book != nil {
		fm.Book.ID = path.Base(path.Dir(pathName))
		fm.Book.contentPath = pathName[strings.LastIndex(pathName, "content/"):strings.LastIndex(pathName, "/")]
		fm.Book.webPath = fm.Book.contentPath[strings.Index(fm.Book.contentPath, "/")+1:]
		log.Println("Found", fm.Book.ID)

		bs.books = append(bs.books, fm.Book)
	}

	return nil
}

func (bs *BookShelf) Books() Books {
	return bs.books
}
