package autodoc

import (
	"context"
	"github.com/peter-mount/documentation/tools/gendoc/util/resource"
	"path"
	"strings"
)

const (
	ResourceManagerKey = "autodoc.ResourceManager"
)

type ResourceManager struct {
	m map[string]resource.Resource
}

func (rm *ResourceManager) Start() error {
	rm.m = make(map[string]resource.Resource)
	return nil
}

func (rm *ResourceManager) GetResources(dir string) resource.Resource {
	var p string
	var pf resource.Resource
	for _, n := range strings.Split(dir, "/") {
		p = path.Join(p, n)
		pf = rm.getResource(pf, p)
	}
	return pf
}

func (rm *ResourceManager) GetResourceIfExists(dir string) (resource.Resource, bool) {
	r, exists := rm.m[dir]
	return r, exists
}

func (rm *ResourceManager) getResource(pf resource.Resource, dir string) resource.Resource {
	if r, exists := rm.m[dir]; exists {
		return r
	}

	r := resource.NewDirectory(path.Base(dir))

	if pf != nil {
		pf.AddChild(r)
	}

	rm.m[dir] = r
	return r
}

func GetResourceManager(ctx context.Context) *ResourceManager {
	if r, ok := ctx.Value(ResourceManagerKey).(*ResourceManager); ok {
		return r
	}
	return nil
}
