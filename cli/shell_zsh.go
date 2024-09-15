package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/alessio/shellescape"
	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func ShellZshInitCmd() command.Command {
	return command.New(
		"Prints the zsh initialization script",
		shellZshInitRun,
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

func shellZshInitRun(cmdCtx command.Context) error {

	printAutoload := cmdCtx.Flags["--print-autoload"].(bool)
	printChpwdHook := cmdCtx.Flags["--print-chpwd-hook"].(bool)
	printExportEnv := cmdCtx.Flags["--print-export-env"].(bool)

	prelude := `
# https://github.com/bbkane/envelope/
#
# To initialize envelope, add this to your configuration (usually ~/.zshrc):
#
# eval "$(envelope shell zsh init)"
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

func ShellZshExportCmd() command.Command {
	return command.New(
		"Print export script",
		shellZshExportRun,
		command.ExistingFlag("--env-name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.Flag(
			"--no-env-no-problem",
			"Exit without an error if the environment doesn't exit. Useful when runnng envelop on chpwd",
			scalar.Bool(
				scalar.Default(false),
			),
			flag.Required(),
		),
	)
}

func shellZshExportRun(cmdCtx command.Context) error {
	return shellZshExportUnexport(cmdCtx, "export")
}

func ShellZshUnexportCmd() command.Command {
	return command.New(
		"Print unexport script",
		shellZshUnexportRun,
		command.ExistingFlag("--env-name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.Flag(
			"--no-env-no-problem",
			"Exit without an error if the environment doesn't exit. Useful when runnng envelop on chpwd",
			scalar.Bool(
				scalar.Default(false),
			),
			flag.Required(),
		),
	)
}

func shellZshUnexportRun(cmdCtx command.Context) error {
	return shellZshExportUnexport(cmdCtx, "unexport")
}

func shellZshExportUnexport(cmdCtx command.Context, scriptType string) error {
	envName := mustGetEnvNameArg(cmdCtx.Flags)
	noEnvNoProblem := cmdCtx.Flags["--no-env-no-problem"].(bool)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	envVars, err := es.EnvVarList(ctx, envName)
	if err != nil {
		if errors.Is(err, domain.ErrEnvNotFound) && noEnvNoProblem {
			return nil
		}
		return fmt.Errorf("could not list env vars: %s: %w", envName, err)
	}

	envRefs, envRefVars, err := es.EnvRefList(ctx, envName)
	if err != nil {
		if errors.Is(err, domain.ErrEnvNotFound) && noEnvNoProblem {
			return nil
		}
		return fmt.Errorf("could not list env refs: %s: %w", envName, err)
	}

	switch scriptType {
	case "export":
		if len(envVars)+len(envRefs) > 0 {
			fmt.Fprintf(cmdCtx.Stdout, "printf '%s:';\n", cmdCtx.AppName)
			for _, ev := range envVars {
				fmt.Fprintf(cmdCtx.Stdout, "printf ' +%s';\n", shellescape.Quote(ev.Name))
				fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(ev.Name), shellescape.Quote(ev.Value))
			}

			for i := range len(envRefs) {
				fmt.Fprintf(cmdCtx.Stdout, "printf ' +%s';\n", shellescape.Quote(envRefs[i].Name))
				fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(envRefs[i].Name), shellescape.Quote(envRefVars[i].Value))
			}
			fmt.Fprintf(cmdCtx.Stdout, "echo;\n")
		}

	case "unexport":
		if len(envVars)+len(envRefs) > 0 {
			fmt.Fprintf(cmdCtx.Stdout, "printf '%s:';\n", cmdCtx.AppName)
			for _, ev := range envVars {
				fmt.Fprintf(cmdCtx.Stdout, "printf ' -%s';\n", shellescape.Quote(ev.Name))
				fmt.Fprintf(cmdCtx.Stdout, "unset %s;\n", shellescape.Quote(ev.Name))
			}

			for _, er := range envRefs {
				fmt.Fprintf(cmdCtx.Stdout, "printf ' -%s';\n", shellescape.Quote(er.Name))
				fmt.Fprintf(cmdCtx.Stdout, "unset %s;\n", shellescape.Quote(er.Name))
			}
			fmt.Fprintf(cmdCtx.Stdout, "echo;\n")
		}
	default:
		return errors.New("unimplemented --script-type: " + scriptType)

	}

	return nil
}