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
		version,
		section.New(
			"Manage Environmental secrets centrally",
			section.CommandMap(warg.VersionCommandMap()),
			section.NewSection(
				"env",
				"Environment commands",
				section.Command("create", cli.EnvCreateCmd()),
				section.Command("delete", cli.EnvDeleteCmd()),
				section.Command("list", cli.EnvListCmd()),
				section.Command("update", cli.EnvUpdateCmd()),
				section.Command("show", cli.EnvShowCmd()),
			),
			section.NewSection(
				"shell",
				"Manipulate the current shell",
				section.NewSection(
					"zsh",
					"Zsh-specific commands",
					section.Command("init", cli.ShellZshInitCmd()),
					section.Command("export", cli.ShellZshExportCmd()),
					section.Command("unexport", cli.ShellZshUnexportCmd()),
				),
			),
			section.NewSection(
				"var",
				"Env vars owned by this environment",
				section.Command("create", cli.VarCreateCmd()),
				section.Command("delete", cli.VarDeleteCmd()),
				section.Command("show", cli.VarShowCmd()),
				section.Command("update", cli.VarUpdateCmd()),
				section.NewSection(
					"ref",
					"Variable References owned by this environment",
					section.Command("create", cli.VarRefCreateCmd()),
					section.Command("delete", cli.VarRefDeleteCmd()),
					section.Command("show", cli.VarRefShowCmd()),
				),
			),
		),
		warg.GlobalFlagMap(warg.ColorFlagMap()),
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
	)
	return &app
}

func main() {
	app := buildApp()
	app.MustRun()
}
