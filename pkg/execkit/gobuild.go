package execkit

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// GoBuild builder
type GoBuild struct {
	Ldflags []string
	Out     string
	Source  string
}

var _ fmt.Stringer = (*GoBuild)(nil)

// BuildVar return ldflag argument for set build variable
func BuildVar(name string, value interface{}) string {
	return fmt.Sprintf("-X %s=%v", name, value)
}

// Command of GoBuild
func (g *GoBuild) Command() *Command {
	return &Command{
		Name:   "go",
		Args:   g.Args(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Args is arguments for go build
func (g *GoBuild) Args() []string {
	args := []string{"build"}
	if len(g.Ldflags) > 0 {
		args = append(args, "-ldflags", strings.Join(g.Ldflags, " "))
	}
	args = append(args, "-o", g.Out, g.Source)
	return args
}

// Run gobuild
func (g *GoBuild) Run(ctx context.Context) error {
	return g.Command().Run(ctx)
}

func (g GoBuild) String() string {
	return g.Command().String()
}