package typical

import (
	"github.com/typical-go/typical-go/examples/provide-constructor/helloworld"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Descriptor of sample
var Descriptor = typcore.Descriptor{
	Name:    "provide-constructor",
	Version: "1.0.0",

	App: typapp.EntryPoint(helloworld.Main),

	BuildTool: typbuildtool.
		BuildSequences(
			typbuildtool.StandardBuild(),
		),
}
