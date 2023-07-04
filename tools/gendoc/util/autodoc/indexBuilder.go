package autodoc

import (
	"context"
	"github.com/peter-mount/documentation/tools/gendoc/util"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"os"
	"path"
	"strings"
	"time"
)

const (
	referenceDir = "reference"
)

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

func GenerateReferenceIndexFile(ctx context.Context, fileName string, fileTime time.Time, title, desc string) error {
	return task.Of(func(ctx context.Context) error {

		return GenerateCustomIndexFile(fileName, fileTime, func(fileName string, fileTime time.Time) error {
			fb := util.ReferenceFileBuilder(title, desc, "manual", 100, fileTime)

			// Get any resources for this location
			if rm := GetResourceManager(ctx); rm != nil {
				if dirResource, exists := rm.GetResourceIfExists(path.Dir(fileName)); exists {
					fb = fb.Then(dirResource.FileBuilder())
				}
			}

			return fb.
				WrapAsFrontMatter().
				FileHandler().
				WriteAlways(fileName, fileTime)
		})

	}).
		WithContext(ctx, ResourceManagerKey).
		QueueWithPriority(90).
		Do(ctx)
}

func GenerateReferenceIndices(fileName string, fileTime time.Time) task.Task {
	return func(ctx context.Context) error {
		// Find location of referenceDir as we only work from there
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

		for i := ri; i < len(p)-1; i++ {
			n := p[i]
			err := GenerateReferenceIndexFile(ctx, path.Join(path.Join(p[:i+1]...), "_index.html"), fileTime, n, "")
			if err != nil {
				return err
			}
		}

		return nil
	}
}
