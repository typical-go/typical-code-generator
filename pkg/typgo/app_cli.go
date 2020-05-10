package typgo

import (
	"os"

	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
	"go.uber.org/dig"
)

func createAppCli(d *Descriptor) *cli.App {
	di := dig.New()
	di.Provide(func() *Descriptor {
		return d
	})

	c := &AppContainer{
		di: di,
	}

	app := cli.NewApp()
	app.Name = d.Name
	app.Usage = "" // NOTE: intentionally blank
	app.Description = d.Description
	app.Before = func(*cli.Context) (err error) {
		if configFile := os.Getenv("CONFIG"); configFile != "" {
			_, err = typcfg.Load(configFile)
		}
		return
	}
	app.Version = d.Version
	app.Action = c.ActionFunc(d.EntryPoint)
	return app
}
