package autodoc

import (
	"context"
	"github.com/peter-mount/documentation/tools/util"
	"github.com/peter-mount/documentation/tools/util/resource"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"os"
	"path"
	"time"
)

func GenerateFileIndexPage(dir, file, asm, buildFileName string, modified time.Time) task.Task {
	fullFileName := path.Join(dir, asm, file, "_index.html")

	return task.Of(func(ctx context.Context) error {
		fi, err := os.Stat(buildFileName)
		if err != nil {
			return err
		}

		rm := GetResourceManager(ctx)

		fileResource := resource.NewFile(fi.Name(), buildFileName[len("content"):], int(fi.Size()))
		rm.GetResources(path.Join(dir, asm)).AddChild(fileResource)

		return GenerateCustomIndexFile(fullFileName, modified, func(indexFileName string, fileTime time.Time) error {
			return util.ReferenceFileBuilder(file, "Generated file for "+asm, "manual", 10, fileTime).
				Then(fileResource.FileBuilder()).
				WrapAsFrontMatter().
				Appendf("{{< book/include src=%q >}}", buildFileName).
				FileHandler().
				WriteAlways(indexFileName, modified)
		})
	}).
		Then(GenerateReferenceIndices(fullFileName, modified))
}
