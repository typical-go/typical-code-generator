package typical

// Autogenerated by Typical-Go. DO NOT EDIT.

import (
	"github.com/typical-go/typical-go/examples/configuration-with-invocation/server"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typgo"
)

func init() {
	typgo.Provide(
		&typgo.Constructor{
			Name: "",
			Fn: func() (cfg *server.Config, err error) {
				cfg = new(server.Config)
				if err = typcfg.Process("SERVER", cfg); err != nil {
					return nil, err
				}
				return
			},
		},
	)
	typgo.Destroy()
}
