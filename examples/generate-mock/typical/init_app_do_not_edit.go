package typical

// Autogenerated by Typical-Go. DO NOT EDIT.

import (
	"github.com/typical-go/typical-go/examples/generate-mock/helloworld"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typdep"
)

func init() {
	typapp.AppendConstructor(
		typdep.NewConstructor(helloworld.NewGreeter),
	)
}
