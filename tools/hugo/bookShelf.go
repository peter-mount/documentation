package hugo

import (
  "github.com/peter-mount/documentation/tools/util/walk"
  "log"
  "os"
  "path"
)

// BookShelf manages all Book's.
// A Book is a _index.html file containing a book: section within it's FrontMatter
type BookShelf struct {
  books Books
}

func (bs *BookShelf) Name() string {
  return "BookShelf"
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
    log.Println("Found", fm.Book.ID)

    bs.books = append(bs.books, fm.Book)
  }

  return nil
}

func (bs *BookShelf) Books() Books {
  return bs.books
}
