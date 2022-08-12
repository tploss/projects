package main

import (
	"fmt"
	"runtime/debug"

	"github.com/alecthomas/kong"
	"gitlab.com/tploss/projects/add"
	"gitlab.com/tploss/projects/edit"
	"gitlab.com/tploss/projects/git"
	"gitlab.com/tploss/projects/globals"
	"gitlab.com/tploss/projects/list"
	"gitlab.com/tploss/projects/man"
	"gitlab.com/tploss/projects/rm"
)

var (
	// version and commitSha should be set via ldflags during build
	version   string
	commitSha string
)

type CLI struct {
	globals.G
	Man man.Man `cmd:"" help:"Generate a manpage for this tool"`

	List   list.ListCmd `cmd:"" aliases:"ls" help:"List projects/repositories (default: list all in yaml format)"`
	Edit   edit.EditCmd `cmd:"" help:"Edit a project"`
	Add    add.AddCmd   `cmd:"" help:"Add a project (alias for edit <project> on non existing project)"`
	Remove rm.RmCmd     `cmd:"" aliases:"rm" help:"Remove project(s) from disk and optionally from config"`
	Git    git.GitCmd   `cmd:"" help:"Run a git command on all repos (of a project)"`
}

func main() {
	vtext := versionText()
	cli := &CLI{}
	ctx := kong.Parse(
		cli,
		kong.Description("Manage your project repositories with ease"),
		kong.UsageOnError(),
		kong.ConfigureHelp(
			kong.HelpOptions{
				Compact: true,
				Summary: false,
			},
		),
		kong.Vars{
			"version": vtext,
		},
		kong.Bind(&cli.G),
	)

	ctx.FatalIfErrorf(ctx.Run())
	ctx.Exit(0)
}

func versionText() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "" {
			version = info.Main.Version
		}
		if commitSha == "" {
			commitSha = getCommitSha(info.Settings)
		}
	}
	return fmt.Sprintf("projects version %s (commit %5s)", version, commitSha)
}

func getCommitSha(bs []debug.BuildSetting) string {
	for _, s := range bs {
		if s.Key == "vcs.revision" {
			return s.Value
		}
	}
	return ""
}
