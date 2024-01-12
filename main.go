package main

import (
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/contained"
	"go.bbkane.com/warg/value/scalar"
)

var version string

func emptyOrNil[T any](iFace interface{}) (T, error) {
	under, ok := iFace.(T)
	if !ok {
		return under, contained.ErrIncompatibleInterface
	}
	return under, nil
}

func datetime() contained.TypeInfo[time.Time] {
	return contained.TypeInfo[time.Time]{
		Description: "datetime in RFC3339 format",
		FromIFace:   emptyOrNil[time.Time],
		FromInstance: func(t time.Time) (time.Time, error) {
			return t, nil
		},
		FromString: func(s string) (time.Time, error) {
			return time.Parse(time.RFC3339, s)
		},
		Empty: func() time.Time {
			return time.Time{}
		},
	}
}

func buildApp() warg.App {

	cwd, err := os.Getwd()
	if err != nil {
		// I don't know when this could happen?
		panic(err)
	}

	envNameFlag := flag.New(
		"Environment name",
		scalar.String(
			scalar.Default(cwd),
		),
		flag.Required(),
	)

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
			scalar.New(
				datetime(),
				scalar.Default(time.Now()),
			),
			flag.Required(),
		),
		"--update-time": flag.New(
			"Update time",
			scalar.New(
				datetime(),
				scalar.Default(time.Now()),
			),
			flag.Required(),
		),
	}

	commonUpdateFlags := flag.FlagMap{
		"--comment": flag.New(
			"Comment",
			scalar.String(),
		),
		"--create-time": flag.New(
			"Create time",
			scalar.New(
				datetime(),
			),
		),
		"--new-name": flag.New(
			"New name",
			scalar.String(
				scalar.Default(cwd),
			),
		),
		"--update-time": flag.New(
			"Update time",
			scalar.New(
				datetime(),
				scalar.Default(time.Now()),
			),
			flag.UnsetSentinel("UNSET"),
		),
	}

	timeoutFlag := flag.FlagMap{
		"--timeout": flag.New(
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(10*time.Minute),
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
					command.ExistingFlag("--name", envNameFlag),
				),
				section.Command(
					"update",
					"Update an environment",
					envUpdateCmd,
					command.ExistingFlags(commonUpdateFlags),
					command.ExistingFlag("--name", envNameFlag),
				),
				section.ExistingFlags(timeoutFlag),
				section.ExistingFlags(sqliteDSN),
				section.Section(
					"var",
					"Environment Variables!",
					section.Command(
						"create",
						"Create an environmental variable",
						envVarCreateCmd,
						command.ExistingFlags(commonCreateFlags),
						command.Flag(
							"--local-value",
							"Value if type is local",
							scalar.String(),
						),
						command.Flag(
							"--name",
							"Env var name",
							scalar.String(),
							flag.Required(),
						),
						command.Flag(
							"--type",
							"Type of env var",
							scalar.String(
								scalar.Choices("local"),
							),
							flag.Required(),
						),
					),
					section.ExistingFlag(
						"--env-name",
						envNameFlag,
					),
				),
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
