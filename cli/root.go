package cli

import (
	"fmt"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func InitZshCmd() command.Command {
	return command.New(
		"Print zsh",
		initZshRun,
		command.Flag(
			"--print-autoload",
			"Include autoload -Uz add-zsh-hook line (might not be needed if you already autoloaded it)",
			scalar.Bool(
				scalar.Default(true),
			),
			flag.Required(),
		),
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

func initZshRun(cmdCtx command.Context) error {

	printAutoload := cmdCtx.Flags["--print-autoload"].(bool)
	printChpwdHook := cmdCtx.Flags["--print-chpwd-hook"].(bool)
	printExportEnv := cmdCtx.Flags["--print-export-env"].(bool)

	prelude := `
# https://github.com/bbkane/envelope/
#
# To initialize envelope, add this to your configuration (usually ~/.zshrc):
#
# eval "$(envelope init)"
#
`
	fmt.Fprint(cmdCtx.Stdout, prelude)

	autoload := `
autoload -Uz add-zsh-hook
`
	if printAutoload {
		fmt.Fprint(cmdCtx.Stdout, autoload)
	}

	chpwdHook := `
add-zsh-hook -Uz chpwd (){
    eval $(envelope env print-script --name "$OLDPWD" --no-env-no-problem true --type unexport)
    eval $(envelope env print-script --name "$PWD" --no-env-no-problem true --type export)
}
`
	if printChpwdHook {
		fmt.Fprint(cmdCtx.Stdout, chpwdHook)
	}

	exportEnv := `
export-env() { eval $(envelope env print-script --name "$1" --no-env-no-problem true --type export) }
unexport-env() { eval $(envelope env print-script --name "$1" --no-env-no-problem true --type unexport) }
`
	if printExportEnv {
		fmt.Fprint(cmdCtx.Stdout, exportEnv)
	}

	return nil
}
