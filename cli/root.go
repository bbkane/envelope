package cli

import (
	"fmt"

	"go.bbkane.com/warg/command"
)

func InitCmd() command.Command {
	return command.New(
		"Print script to run envelope when the directory changes. Only zsh supported (for now)",
		initRun,
	)
}

func initRun(cmdCtx command.Context) error {
	script := `
# https://github.com/bbkane/envelope/
#
# To initialize envelope, add this to your configuration (usually ~/.zshrc):
#
# eval "$(envelope init)"
#
autoload -Uz add-zsh-hook
add-zsh-hook -Uz chpwd (){
    eval $(envelope env print-script --name "$PWD" --no-env-no-problem true --type export)
    eval $(envelope env print-script --name "$OLDPWD" --no-env-no-problem true --type unexport)
}
`
	fmt.Fprintln(cmdCtx.Stdout, script)
	return nil
}
