package typgo

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typannot"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typtmpl"
	"github.com/typical-go/typical-go/pkg/typvar"
	"github.com/urfave/cli/v2"
)

var (
	_ typcore.AppLauncher       = (*Descriptor)(nil)
	_ typcore.BuildToolLauncher = (*Descriptor)(nil)

	_ Utility        = (*Descriptor)(nil)
	_ Preconditioner = (*Descriptor)(nil)
)

type (

	// Descriptor describe the project
	Descriptor struct {

		// Name of the project (OPTIONAL).
		// It should be a characters with/without underscore or dash.
		// By default, project name is same with project folder
		Name string

		// Description of the project (OPTIONAL).
		Description string

		// Version of the project (OPTIONAL).
		// By default it is 0.0.1
		Version string

		BuildSequences []interface{}

		Utility Utility

		Layouts []string

		SkipPrecond bool

		EntryPoint interface{}

		Configurer typcfg.Configurer
	}
)

// LaunchApp to launch the app
func (d *Descriptor) LaunchApp() (err error) {
	if err = d.Validate(); err != nil {
		return
	}
	return createAppCli(d).Run(os.Args)
}

// LaunchBuildTool to launch the build tool
func (d *Descriptor) LaunchBuildTool() (err error) {
	if err := d.Validate(); err != nil {
		return err
	}

	return createBuildToolCli(d).Run(os.Args)
}

// Validate context
func (d *Descriptor) Validate() (err error) {
	if d.Version == "" {
		d.Version = "0.0.1"
	}

	if !ValidateName(d.Name) {
		return errors.New("Descriptor: bad name")
	}

	if len(d.BuildSequences) < 1 {
		return errors.New("Descriptor:No build-sequence")
	}

	for _, module := range d.BuildSequences {
		if err = common.Validate(module); err != nil {
			return err
		}
	}

	return
}

// ValidateName to validate valid descriptor name
func ValidateName(name string) bool {
	if name == "" {
		return false
	}

	r, _ := regexp.Compile(`^[a-zA-Z\_\-]+$`)
	if !r.MatchString(name) {
		return false
	}
	return true
}

// Commands to return command
func (d *Descriptor) Commands(c *BuildTool) (cmds []*cli.Command) {
	cmds = []*cli.Command{
		cmdTest(c),
		cmdRun(c),
		cmdPublish(c),
		cmdClean(c),
	}

	if d.Utility != nil {
		for _, cmd := range d.Utility.Commands(c) {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

// Precondition for this project
func (d *Descriptor) Precondition(c *PrecondContext) (err error) {
	if d.SkipPrecond {
		c.Info("Skip the preconditon")
		return
	}

	if d.Configurer != nil {
		if err = typcfg.Write(typvar.ConfigFile, d.Configurer); err != nil {
			return
		}
	}

	appPrecond := d.appPrecond(c)
	if appPrecond.NotEmpty() {
		c.AppendTemplate(appPrecond)
	}

	typcfg.Load(typvar.ConfigFile)

	return
}

func (d *Descriptor) appPrecond(c *PrecondContext) *typtmpl.AppPrecond {
	var (
		ctors    []*typtmpl.Ctor
		cfgCtors []*typtmpl.CfgCtor
		dtors    []*typtmpl.Dtor
	)

	store := c.ASTStore()

	ctorAnnots, errs := typannot.GetCtors(store)
	for _, a := range ctorAnnots {
		ctors = append(ctors, &typtmpl.Ctor{
			Name: a.Name,
			Def:  fmt.Sprintf("%s.%s", a.Decl.Pkg, a.Decl.Name),
		})
	}

	dtorAnnots, errs := typannot.GetDtors(store)
	for _, a := range dtorAnnots {
		dtors = append(dtors, &typtmpl.Dtor{
			Def: fmt.Sprintf("%s.%s", a.Decl.Pkg, a.Decl.Name),
		})
	}

	for _, err := range errs {
		c.Warnf("App-Precond: %s", err.Error())
	}

	if d.Configurer != nil {
		for _, cfg := range d.Configurer.Configurations() {
			specType := reflect.TypeOf(cfg.Spec).String()
			cfgCtors = append(cfgCtors, &typtmpl.CfgCtor{
				Name:      cfg.CtorName,
				Prefix:    cfg.Name,
				SpecType:  specType,
				SpecType2: specType[1:],
			})
		}
	}

	return &typtmpl.AppPrecond{
		Ctors:    ctors,
		CfgCtors: cfgCtors,
		Dtors:    dtors,
	}
}