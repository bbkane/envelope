package main

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

var version string

func buildApp() warg.App {

	cwd, err := os.Getwd()
	if err != nil {
		// I don't know when this could happen?
		panic(err)
	}

	dbPath, err := homedir.Expand("~/.config/namedenv.db")
	if err != nil {
		panic(err)
	}

	sqliteFlags := flag.FlagMap{
		"--sqlite-dsn": flag.New(
			"Sqlite DSN. Usually the file name",
			scalar.String(
				scalar.Default(dbPath),
			),
			flag.Required(),
		),
	}

	app := warg.New(
		"namedenv",
		section.New(
			"Manage Environmental secrets centrally",
			section.ExistingCommand("version", warg.VersionCommand()),
			section.ExistingFlag("--color", warg.ColorFlag()),
			section.Section(
				"env",
				"Environment commands",
				section.Command(
					"create",
					"Create an environment",
					envCreateCmd,
					command.Flag(
						"--name",
						"Name of environment",
						scalar.String(
							scalar.Default(cwd),
						),
						flag.Alias("-n"),
						flag.Required(),
					),
				),
				section.ExistingFlags(sqliteFlags),
			),
		),
		warg.OverrideVersion(version),
	)
	return app
}

func main() {
	app := buildApp()
	app.MustRun()
}
