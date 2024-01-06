package main

import (
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/section"
)

var version string

func buildApp() warg.App {
	app := warg.New(
		"example-go-cli",
		section.New(
			"Example Go CLI",
			section.Command(
				"hello",
				"Say hello",
				hello,
			),
		),
		warg.AddColorFlag(),
		warg.AddVersionCommand(version),
	)
	return app
}

func main() {
	app := buildApp()
	app.MustRun(os.Args, os.LookupEnv)
}
