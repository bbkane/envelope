package main

import (
	"go.bbkane.com/namedenv/cli"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/section"
)

var version string

func buildApp() *warg.App {

	app := warg.New(
		"namedenv",
		section.New(
			"Manage Environmental secrets centrally",
			section.ExistingCommand("version", warg.VersionCommand()),
			section.ExistingFlag("--color", warg.ColorFlag()),
			section.Section(
				"env",
				"Environment commands",
				section.ExistingCommand(
					"create",
					cli.EnvCreateCmd(),
				),
				section.ExistingCommand(
					"delete",
					cli.EnvDeleteCmd(),
				),
				section.ExistingCommand(
					"update",
					cli.EnvUpdateCmd(),
				),
				section.ExistingCommand(
					"print-script",
					cli.EnvPrintScriptCmd(),
				),
				section.ExistingCommand(
					"show",
					cli.EnvShowCmd(),
				),
				section.Section(
					"localvar",
					"Env vars owned by this environment",
					section.ExistingCommand(
						"create",
						cli.EnvLocalVarCreateCmd(),
					),
					section.ExistingCommand(
						"show",
						cli.EnvLocalVarShowCmd(),
					),
				),
			),
			section.Section(
				"keyring",
				"Work with the OS Keyring",
				section.ExistingCommand(
					"create",
					cli.KeyringCreateCmd(),
				),
			),
		),
		warg.OverrideVersion(version),
	)
	return &app
}

func main() {
	app := buildApp()
	app.MustRun()
}
