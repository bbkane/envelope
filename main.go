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

	sqliteDSN := flag.FlagMap{
		"--sqlite-dsn": flag.New(
			"Sqlite DSN. Usually the file name",
			scalar.String(
				scalar.Default(dbPath),
			),
			flag.Required(),
		),
	}

	// most tables have these, so let's just re-use the definition
	commonCreateFlags := flag.FlagMap{
		"--comment": flag.New(
			"Comment",
			scalar.String(),
		),
		"--create-time": flag.New(
			"Create time",
			scalar.String(
				scalar.Default("bob"), // TODO: make current time
			),
			flag.Required(),
		),
		"--update-time": flag.New(
			"Update time",
			scalar.String(
				scalar.Default("bob"), // TODO: make current time
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
					command.ExistingFlags(commonCreateFlags),
					command.Flag(
						"--name",
						"Environment name",
						scalar.String(
							scalar.Default(cwd),
						),
						flag.Alias("-n"),
						flag.Required(),
					),
				),
				section.ExistingFlags(sqliteDSN),
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
