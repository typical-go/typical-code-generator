package typmock

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/typical-go/typical-go/pkg/typbuild/prebld"
)

type mockTarget struct {
	srcPkg  string
	srcName string
	mockPkg string
	dest    string
}

func createMockTarget(c *Context, decl *prebld.Declaration) *mockTarget {
	var (
		pkg     = decl.File.Name.Name
		dir     = filepath.Dir(decl.Path)
		dirDest = dir[:len(dir)-len(pkg)]
		srcPkg  = fmt.Sprintf("%s/%s", c.ModulePackage, dir)
		mockPkg = fmt.Sprintf("%s_%s", c.MockFolder, pkg)
		dest    = fmt.Sprintf("%s%s/%s.go", dirDest, mockPkg, strcase.ToSnake(decl.SourceName))
	)
	return &mockTarget{
		srcPkg:  srcPkg,
		srcName: decl.SourceName,
		mockPkg: mockPkg,
		dest:    dest,
	}
}
