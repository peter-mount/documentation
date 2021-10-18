package autodoc

import (
  "context"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/task"
  "os"
  "path"
  "time"
)

func GenerateIndexPage(dir, file, asm, buildFileName string, modified time.Time) task.Task {
  return func(ctx context.Context) error {
    fi, err := os.Stat(buildFileName)
    if err != nil {
      return err
    }

    rm := GetResourceManager(ctx)

    res := rm.GetResources(dir).
      Add(fi.Name(), buildFileName[len("content"):], int(fi.Size()))

    return util.GenerateCustomIndexFile(path.Join(dir, asm, file, "_index.html"), modified, func(indexFileName string, fileTime time.Time) error {
      return util.ReferenceFileBuilder(file, "Generated file for "+asm, "manual", 10).
        Then(res.FileBuilder()).
        WrapAsFrontMatter().
        Appendf("{{< book/include src=%q >}}", buildFileName).
        FileHandler().
        WriteAlways(indexFileName, modified)
    })
  }
}
