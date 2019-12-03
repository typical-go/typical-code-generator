package typenv

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	Layout = struct {
		App      string
		Bin      string
		Cmd      string
		Metadata string
		Mock     string
		Release  string
	}{
		App:      "app",
		Cmd:      "cmd",
		Bin:      "bin",
		Metadata: ".typical-metadata",
		Mock:     "mock",
		Release:  "release",
	}

	Readme = "README.md"
	Name   = name()

	AppBin      = fmt.Sprintf("%s/%s", Layout.Bin, Name)
	AppMainPath = fmt.Sprintf("%s/%s", Layout.Cmd, Name)

	BuildToolBin      = fmt.Sprintf("%s/%s-buildtool", Layout.Bin, Name)
	BuildToolMainPath = fmt.Sprintf("%s/%s-buildtool", Layout.Cmd, Name)

	PrebuilderBin      = fmt.Sprintf("%s/%s-prebuilder", Layout.Bin, Name)
	PrebuilderMainPath = fmt.Sprintf("%s/%s-prebuilder", Layout.Cmd, Name)

	Dependency     = "dependency"
	DependencyPath = fmt.Sprintf("internal/%s", Dependency)

	ContextFile  = "typical/context.go"
	ChecksumFile = Layout.Metadata + "/checksum"
)

func name() (s string) {
	var err error
	if s, err = os.Getwd(); err != nil {
		return "noname"
	}
	return filepath.Base(s)
}
