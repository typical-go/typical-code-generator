package typprebuilder

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/utility/coll"
	"github.com/typical-go/typical-go/pkg/utility/filekit"

	"github.com/typical-go/typical-go/pkg/typenv"
	"github.com/typical-go/typical-go/pkg/utility/debugkit"
	"github.com/typical-go/typical-go/pkg/utility/golang"
)

type mockTarget struct {
	ApplicationImports coll.KeyStrings
	MockTargets        []string
}

func (g mockTarget) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate mock target")()
	src := golang.NewSource(typenv.Dependency)
	src.Imports = g.ApplicationImports
	for _, mockTarget := range g.MockTargets {
		src.Init.Append(fmt.Sprintf("typical.Context.MockTargets.Append(\"%s\")", mockTarget))
	}
	if err = filekit.Write(target, src); err != nil {
		return
	}
	return goimports(target)
}
