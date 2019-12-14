package typprebuilder

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typctx"
	"github.com/typical-go/typical-go/pkg/typenv"
	"github.com/typical-go/typical-go/pkg/typprebuilder/metadata"
	"github.com/typical-go/typical-go/pkg/typprebuilder/walker"
	"github.com/typical-go/typical-go/pkg/utility/coll"
	"github.com/typical-go/typical-go/pkg/utility/filekit"
)

// Run the prebuilder
func Run(ctx *typctx.Context) {
	var err error
	var buildCmds []string
	var filenames coll.Strings
	var autowires Autowires
	var automocks Automocks
	var configuration bool
	var buildCommands bool
	readmeFile := !filekit.IsExist(typenv.Readme)
	if filenames, err = projectFiles(typenv.Layout.App); err != nil {
		return
	}
	walker := walker.New(filenames)
	walker.AddFuncDeclListener(&autowires)
	walker.AddTypeSpecListener(&automocks)
	if err = walker.Walk(); err != nil {
		return
	}
	for _, cmd := range typbuildtool.BuildCommands(ctx) {
		for _, subcmd := range cmd.Subcommands {
			buildCmds = append(buildCmds, fmt.Sprintf("%s_%s", cmd.Name, subcmd.Name))
		}
	}
	if buildCommands, err = metadata.Update("build_commands", buildCmds); err != nil {
		log.Fatal(err.Error())
	}
	if _, err = metadata.Update("mock_target", automocks); err != nil {
		log.Fatal(err.Error())
	}
	if _, err = Generate("constructor", constructor{ProjectPackage: ctx.Package, Constructors: autowires}); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Build the build-tool")
	cmd := exec.Command("go", "build",
		"-o", typenv.BuildToolBin,
		"./"+typenv.BuildToolMainPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err.Error())
	}

	if ctx.ReadmeGenerator != nil && (buildCommands || configuration || readmeFile) {
		log.Info("Generate readme")
		cmd := exec.Command(typenv.BuildToolBin, "readme")
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func contextChecksum() bool {
	// NOTE: context checksum is passed by typicalw
	if len(os.Args) > 1 {
		return os.Args[1] == "1"
	}
	return false
}

func goimports(filename string) error {
	cmd := exec.Command(fmt.Sprintf("%s/bin/goimports", build.Default.GOPATH),
		"-w", filename)
	return cmd.Run()
}
