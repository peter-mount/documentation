package util

import (
  "fmt"
  "os"
  "path"
  "strings"
  "time"
)

func GenerateReferenceIndex(fileName string, fileTime time.Time) error {
  if i := strings.LastIndex(fileName, referenceDir); i > 0 {
    p := fileName[:i+len(referenceDir)]
    if fi, err := os.Stat(p); err != nil {
      if !os.IsNotExist(err) {
        return err
      }
    } else if !fi.IsDir() {
      return fmt.Errorf("%s exists & is a file", fileName[:i+len(referenceDir)])
    }

    // Create our default reference/_index.html
    return GenerateReferenceIndexFile(path.Join(p, "_index.html"), fileTime, "Reference", "Reference", "Generated tables & files from the documentation")
  }

  return nil
}

type CustomIndexFileGenerator func(fileName string, fileTime time.Time) error

func GenerateCustomIndexFile(fileName string, fileTime time.Time, f CustomIndexFileGenerator) error {
  if _, err := os.Stat(fileName); err != nil {
    if !os.IsNotExist(err) {
      return err
    }
    return f(fileName, fileTime)
  }
  return nil
}

func GenerateReferenceIndexFile(fileName string, fileTime time.Time, title, linkTitle, desc string) error {
  return GenerateCustomIndexFile(fileName, fileTime, func(fileName string, fileTime time.Time) error {
    return StringFileHandler(fmt.Sprintf(
      "---\n# This file is automatically generated\ntype: \"manual\"\ntitle: %q\nlinkTitle: %q\ndescription: %q\n---",
      title, linkTitle, desc,
    )).
      writeAlways(fileName, fileTime)
  })
}

func GenerateReferenceIndices(fileName string, fileTime time.Time) error {
  if err := GenerateReferenceIndex(fileName, fileTime); err != nil {
    return err
  }

  p := strings.Split(fileName, "/")
  ri := -1
  for i, e := range p {
    if e == referenceDir {
      ri = i
    }
  }
  if ri == -1 {
    return nil
  }

  for i := ri + 1; i < len(p)-1; i++ {
    n := p[i]
    err := GenerateReferenceIndexFile(path.Join(path.Join(p[:i+1]...), "_index.html"), fileTime, n, n, "")
    if err != nil {
      return err
    }
  }

  return nil
}