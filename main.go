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
			section.ExistingFlag("--color", warg.ColorFlag()),
			section.ExistingCommand("init", cli.InitCmd()),
			section.Section(
				"env",
				"Environment commands",
				section.ExistingCommand("create", cli.EnvCreateCmd()),
				section.ExistingCommand("delete", cli.EnvDeleteCmd()),
				section.ExistingCommand("list", cli.EnvListCmd()),
				section.ExistingCommand("update", cli.EnvUpdateCmd()),
				section.ExistingCommand("print-script", cli.EnvPrintScriptCmd()),
				section.ExistingCommand("show", cli.EnvShowCmd()),
				section.Section(
					"localvar",
					"Env vars owned by this environment",
					section.ExistingCommand("create", cli.EnvLocalVarCreateCmd()),
					section.ExistingCommand("delete", cli.EnvLocalVarDeleteCmd()),
					section.ExistingCommand("show", cli.EnvLocalVarShowCmd()),
				),
			),
			section.Section(
				"keyring",
				"Work with the OS Keyring",
				section.ExistingCommand("create", cli.KeyringCreateCmd()),
				section.ExistingCommand("list", cli.KeyringListCmd()),
			),
		),
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
