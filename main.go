package main

import (
	"go.bbkane.com/envelope/cli"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/help/allcommands"
	"go.bbkane.com/warg/help/detailed"
	"go.bbkane.com/warg/section"
)

var version string

func buildApp() *warg.App {

	app := warg.New(
		"envelope",
		section.New(
			"Manage Environmental secrets centrally",
			section.ExistingCommand("version", warg.VersionCommand()),
			section.Section(
				"env",
				"Environment commands",
				section.ExistingCommand("create", cli.EnvCreateCmd()),
				section.ExistingCommand("delete", cli.EnvDeleteCmd()),
				section.ExistingCommand("list", cli.EnvListCmd()),
				section.ExistingCommand("update", cli.EnvUpdateCmd()),
				section.ExistingCommand("show", cli.EnvShowCmd()),
			),
			section.Section(
				"shell",
				"Manipulate the current shell",
				section.Section(
					"zsh",
					"Zsh-specific commands",
					section.ExistingCommand("init", cli.ShellZshInitCmd()),
					section.ExistingCommand("export", cli.ShellZshExportCmd()),
					section.ExistingCommand("unexport", cli.ShellZshUnexportCmd()),
				),
			),
			section.Section(
				"var",
				"Env vars owned by this environment",
				section.ExistingCommand("create", cli.VarCreateCmd()),
				section.ExistingCommand("delete", cli.VarDeleteCmd()),
				section.ExistingCommand("show", cli.VarShowCmd()),
				section.ExistingCommand("update", cli.VarUpdateCmd()),
				section.Section(
					"ref",
					"Variable References owned by this environment",
					section.ExistingCommand("create", cli.VarRefCreateCmd()),
					section.ExistingCommand("delete", cli.VarRefDeleteCmd()),
					section.ExistingCommand("show", cli.VarRefShowCmd()),
				),
			),
		),
		warg.ExistingGlobalFlag("--color", warg.ColorFlag()),
		warg.OverrideHelpFlag(
			[]help.HelpFlagMapping{
				{Name: "detailed", CommandHelp: detailed.DetailedCommandHelp, SectionHelp: detailed.DetailedSectionHelp},
				{Name: "outline", CommandHelp: help.OutlineCommandHelp, SectionHelp: help.OutlineSectionHelp},
				// allcommands list child commands, so it doesn't really make sense for a command
				{Name: "allcommands", CommandHelp: detailed.DetailedCommandHelp, SectionHelp: allcommands.AllCommandsSectionHelp},
			},
			"detailed",
			"--help",
			"Print help",
			flag.Alias("-h"),
		),
		warg.OverrideVersion(version),
	)
	return &app
}

func main() {
	app := buildApp()
	app.MustRun()
}
