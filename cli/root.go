package cli

import (
	"fmt"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func InitCmd() command.Command {
	return command.New(
		"Print script to run envelope when the directory changes. Only zsh supported (for now)",
		initRun,
		command.Flag(
			"--print-export-env",
			"Include export-env/unexport-env to easily use envs from the CLI",
			scalar.Bool(
				scalar.Default(true),
			),
			flag.Required(),
		),
		command.Flag(
			"--print-chpwd-hook",
			"Include hook to export/unexport envs when changing directories",
			scalar.Bool(
				scalar.Default(true),
			),
			flag.Required(),
		),
	)
}

func initRun(cmdCtx command.Context) error {

	printChpwdHook := cmdCtx.Flags["--print-chpwd-hook"].(bool)
	printExportEnv := cmdCtx.Flags["--print-export-env"].(bool)

	chpwdHook := `
# https://github.com/bbkane/envelope/
#
# To initialize envelope, add this to your configuration (usually ~/.zshrc):
#
# eval "$(envelope init)"
#
autoload -Uz add-zsh-hook
add-zsh-hook -Uz chpwd (){
    eval $(envelope env print-script --name "$OLDPWD" --no-env-no-problem true --type unexport)
    eval $(envelope env print-script --name "$PWD" --no-env-no-problem true --type export)
}
`
	if printChpwdHook {
		fmt.Fprintln(cmdCtx.Stdout, chpwdHook)
	}

	exportEnv := `
export-env() { eval $(envelope env print-script --name "$1" --no-env-no-problem true --type export) }
unexport-env() { eval $(envelope env print-script --name "$1" --no-env-no-problem true --type unexport) }
`
	if printExportEnv {
		fmt.Fprintln(cmdCtx.Stdout, exportEnv)
	}

	return nil
}
