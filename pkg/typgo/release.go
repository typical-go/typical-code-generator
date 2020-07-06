package typgo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/typical-go/typical-go/pkg/git"
	"github.com/urfave/cli/v2"
)

type (
	// Releaser responsible to release
	Releaser interface {
		Release(*ReleaseContext) error
	}

	// ReleaseContext contain data for release
	ReleaseContext struct {
		*Context

		Alpha   bool
		Tag     string
		GitLogs []*git.Log
	}

	// Releases for composite release
	Releases []Releaser

	// ReleaseFn release function
	ReleaseFn func(*ReleaseContext) error

	releaserImpl struct {
		fn ReleaseFn
	}
)

var _ Releaser = (Releases)(nil)

//
// releaserImpl
//

// NewRelease return new instance of Releaser
func NewRelease(fn ReleaseFn) Releaser {
	return &releaserImpl{fn: fn}
}

func (r *releaserImpl) Release(c *ReleaseContext) error {
	return r.fn(c)
}

//
// Releaser
//

// Release the releasers
func (r Releases) Release(c *ReleaseContext) (err error) {
	for _, releaser := range r {
		if err = releaser.Release(c); err != nil {
			return
		}
	}
	return
}

//
// Command
//

func cmdRelease(c *BuildCli) *cli.Command {
	return &cli.Command{
		Name:  "release",
		Usage: "Release the project",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "skip-test", Usage: "skip test"},
			&cli.BoolFlag{Name: "skip-compile", Usage: "skip compile"},
			&cli.BoolFlag{Name: "force", Usage: "Release by passed all validation"},
			&cli.BoolFlag{Name: "alpha", Usage: "Release for alpha version"},
		},
		Action: c.ActionFn("RELEASE", release),
	}
}

func release(c *Context) (err error) {
	if c.Release == nil {
		return errors.New("No Releaser")
	}

	if !c.Bool("skip-test") {
		if err = test(c); err != nil {
			return
		}
	}

	if !c.Bool("skip-compile") {
		if err = compile(c); err != nil {
			return
		}
	}

	ctx := c.Ctx()

	if err = git.Fetch(ctx); err != nil {
		return fmt.Errorf("Failed git fetch: %w", err)
	}
	defer git.Fetch(ctx)

	force := c.Bool("force")
	alpha := c.Bool("alpha")
	tag := releaseTag(c, alpha)

	status := git.Status(ctx)
	if status != "" && !force {
		return fmt.Errorf("Please commit changes first:\n%s", status)
	}

	latest := git.LatestTag(ctx)
	if latest == tag && !force {
		return fmt.Errorf("%s already released", latest)
	}

	gitLogs := git.RetrieveLogs(ctx, latest)
	if len(gitLogs) < 1 && !force {
		return errors.New("No change to be released")
	}

	return c.Release.Release(&ReleaseContext{
		Context: c,
		Alpha:   alpha,
		Tag:     tag,
		GitLogs: gitLogs,
	})

}

func releaseTag(c *Context, alpha bool) string {
	version := "0.0.1"
	if c.Descriptor.Version != "" {
		version = c.Descriptor.Version
	}

	var builder strings.Builder
	builder.WriteString("v")
	builder.WriteString(version)
	// if c.BuildTool.IncludeBranch {
	// 	builder.WriteString("_")
	// 	builder.WriteString(git.Branch(c.Context))
	// }
	// if c.BuildTool.IncludeCommitID {
	// 	builder.WriteString("_")
	// 	builder.WriteString(git.LatestCommit(c.Context))
	// }
	if alpha {
		builder.WriteString("_alpha")
	}
	return builder.String()
}