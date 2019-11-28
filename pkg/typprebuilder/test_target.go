package typprebuilder

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typenv"
	"github.com/typical-go/typical-go/pkg/utility/bash"
	"github.com/typical-go/typical-go/pkg/utility/debugkit"
	"github.com/typical-go/typical-go/pkg/utility/filekit"
	"github.com/typical-go/typical-go/pkg/utility/golang"
)

type testTarget struct {
	ContextImport string
	Packages      []string
}

func (g testTarget) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate test target")()
	src := golang.NewSource(typenv.Dependency)
	src.Imports.Add("", g.ContextImport)
	for _, pkg := range g.Packages {
		src.Init.Append(fmt.Sprintf("typical.Context.TestTargets.Append(\"./%s\")", pkg))
	}
	if err = filekit.Write(target, src); err != nil {
		return
	}
	return bash.GoImports(target)
}
