package typbuild

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typbuild/prebld"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typdep"
	"github.com/urfave/cli/v2"
)

// Context of build
type Context struct {
	*typcore.TypicalContext
	Declarations []*prebld.Declaration
}

// DeclFunc to handle declaration
type DeclFunc func(*prebld.Declaration) error

// AnnotationFunc to handle annotation
type AnnotationFunc func(decl *prebld.Declaration, ann *prebld.Annotation) error

// EachDecl to handle each declaration
func (c *Context) EachDecl(fn DeclFunc) (err error) {
	for _, decl := range c.Declarations {
		if err = fn(decl); err != nil {
			return
		}
	}
	return
}

// EachAnnotation to handle each annotation
func (c *Context) EachAnnotation(name string, declType prebld.DeclType, fn AnnotationFunc) (err error) {
	return c.EachDecl(func(decl *prebld.Declaration) (err error) {
		annotation := decl.Annotations.Get(name)
		if annotation != nil {
			if decl.Type == declType {
				if err = fn(decl, annotation); err != nil {
					return
				}
			} else {
				log.Warnf("[%s] has no effect to %s:%s", name, declType, decl.SourceName)
			}
		}
		return
	})
}

// ActionFunc to return ActionFunc to invoke function fn
func (c *Context) ActionFunc(v interface{}) func(*cli.Context) error {
	return func(cliCtx *cli.Context) (err error) {
		if invocation, ok := v.(*typdep.Invocation); ok {
			return c.Invoke(cliCtx, invocation)
		}
		return c.Invoke(cliCtx, typdep.NewInvocation(v))
	}
}

// Invoke function
func (c *Context) Invoke(cliCtx *cli.Context, invocation *typdep.Invocation) (err error) {
	di := typdep.New()

	if err = typdep.NewConstructor(func() *cli.Context {
		return cliCtx
	}).Provide(di); err != nil {
		return
	}

	// provide functions
	if c.Configuration != nil {
		if err = typdep.ProvideAll(di, c.Configuration.Store().Provide()...); err != nil {
			return
		}
	}

	startFn := func() error { return invocation.Invoke(di) }
	for _, err := range common.StartGracefully(startFn, nil) {
		log.Error(err.Error())
	}
	return
}
