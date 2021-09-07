package walk

import (
  "errors"
  "os"
  "path/filepath"
  "strings"
)

type PathWalker func(path string, info os.FileInfo) error

type PathPredicate func(path string, info os.FileInfo) bool

func (a PathPredicate) Not() PathPredicate {
  return func(path string, info os.FileInfo) bool {
    return !a(path, info)
  }
}

func NewPathWalker() PathWalker {
  return func(path string, info os.FileInfo) error {
    return nil
  }
}

func (a PathWalker) Then(b PathWalker) PathWalker {
  return func(path string, info os.FileInfo) error {
    err := a(path, info)
    if err != nil {
      return err
    }
    return b(path, info)
  }
}

var (
  // Used to stop the chain when predicate is false
  predicateFail = errors.New("predicate fail")
)

func (a PathWalker) If(p PathPredicate) PathWalker {
  return a.Then(func(path string, info os.FileInfo) error {
    if p(path, info) {
      return a(path, info)
    }
    return predicateFail
  })
}

func (a PathWalker) IfNot(p PathPredicate) PathWalker {
  return a.If(p.Not())
}

func (a PathWalker) IsDir() PathWalker {
  return a.If(func(_ string, info os.FileInfo) bool {
    return info.IsDir()
  })
}

func (a PathWalker) IsFile() PathWalker {
  return a.If(func(_ string, info os.FileInfo) bool {
    return !info.IsDir()
  })
}

func (a PathWalker) PathContains(s string) PathWalker {
  return a.If(func(path string, _ os.FileInfo) bool {
    return strings.Contains(path, s)
  })
}

func (a PathWalker) PathNotContain(s string) PathWalker {
  return a.If(func(path string, _ os.FileInfo) bool {
    return !strings.Contains(path, s)
  })
}

func (a PathWalker) PathHasSuffix(s string) PathWalker {
  return a.If(func(path string, _ os.FileInfo) bool {
    return strings.HasSuffix(path, s)
  })
}

func (a PathWalker) PathHasNotSuffix(s string) PathWalker {
  return a.If(func(path string, _ os.FileInfo) bool {
    return !strings.HasSuffix(path, s)
  })
}

func (a PathWalker) Walk(root string) error {
  return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
    // Any error walking to the file or directory will abort the walk immediately
    if err != nil {
      return err
    }

    // Enter our PathWalker chain
    err= a(path, info)
    // Don't pass this to the Walker as it just means we have stopped processing
    if err == predicateFail {
      return nil
    }
    return err
  })
}
