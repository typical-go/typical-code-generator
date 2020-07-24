package typgo

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/typical-go/typical-go/pkg/common"
)

type (
	// Descriptor describe the project
	Descriptor struct {
		// Name of the project (OPTIONAL). It should be a characters with/without underscore or dash.
		// By default, project name is same with project folder
		Name string
		// Version of the project (OPTIONAL). By default it is 0.0.1
		Version string
		Layouts []string
		Cmds    []Cmd
	}
)

// Start typical build-tool
func Start(d *Descriptor) {
	if envmap, _ := common.CreateEnvMapFromFile(".env"); envmap != nil {
		if err := envmap.Setenv(); err == nil {
			printEnv(os.Stdout, envmap)
		}
	}

	if err := createBuildSys(d).app().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func printEnv(w io.Writer, envs map[string]string) {
	color.New(color.FgGreen).Fprint(w, "ENV")
	fmt.Fprint(w, ": ")

	for key := range envs {
		fmt.Fprintf(w, "+%s ", key)
	}
	fmt.Fprintln(w)
}