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
					"update",
					cli.EnvUpdateCmd(),
				),
				section.ExistingCommand(
					"show",
					cli.EnvShowCmd(),
				),
				section.Section(
					"print-script",
					"Print a script from a stored environment",
					section.ExistingCommand(
						"export",
						cli.EnvPrintScriptExportCmd(),
					),
				),

				section.Section(
					"var",
					"Environment Variables!",
					section.Section(
						"create",
						"Create an environmetnal variable",
						section.ExistingCommand(
							"local",
							cli.EnvVarCreateLocalCmd(),
						),
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
