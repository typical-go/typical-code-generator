package typapp

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/typical-go/typical-go/pkg/exor"
	"github.com/typical-go/typical-go/pkg/typapp/internal/tmpl"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Precondition the app
func (a *TypicalApp) Precondition(c *typbuildtool.Context) (err error) {
	var constructors []string

	if err = c.Ast().EachAnnotation("constructor", typast.FunctionType, func(decl *typast.Declaration, ann *typast.Annotation) (err error) {
		constructors = append(constructors, fmt.Sprintf("%s.%s", decl.File.Name, decl.SourceName))
		return
	}); err != nil {
		return
	}

	if c.ConfigManager != nil {
		for _, bean := range c.Configurations() {
			constructors = append(constructors, configDefinition(bean))
		}
	}

	c.Info("Generate constructors")
	target := "typical/init_app_do_not_edit.go"
	if err = a.generateConstructor(c, target, constructors); err != nil {
		return
	}
	return
}

func configDefinition(bean *typcore.Configuration) string {
	typ := reflect.TypeOf(bean.Spec()).String()
	return fmt.Sprintf(`func(cfgMngr typcore.ConfigManager) (%s, error){
		cfg, err := cfgMngr.RetrieveConfig("%s")
		if err != nil {
			return nil, err
		}
		return  cfg.(%s), nil 
	}`, typ, bean.Name(), typ)
}

func (a *TypicalApp) generateConstructor(c *typbuildtool.Context, target string, constructors []string) (err error) {
	ctx := context.Background()
	imports := []string{}
	for _, dir := range c.ProjectDirs {
		if !strings.Contains(dir, "internal") {
			imports = append(imports, fmt.Sprintf("%s/%s", c.ProjectPackage, dir))
		}
	}
	if err = exor.NewWriteTemplate(target, tmpl.Constructor, tmpl.ConstructorData{
		Imports:      imports,
		Constructors: constructors,
	}).Execute(ctx); err != nil {
		return
	}
	cmd := exec.CommandContext(ctx,
		fmt.Sprintf("%s/bin/goimports", build.Default.GOPATH),
		"-w", target)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
