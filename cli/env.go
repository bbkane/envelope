package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alessio/shellescape"
	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/domain"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvCreateCmd() command.Command {
	return command.New(
		"Create an environment",
		envCreateRun,
		command.ExistingFlags(commonCreateFlagMap()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
	)
}

func envCreateRun(cmdCtx command.Context) error {

	commonCreateArgs := mustGetCommonCreateArgs(cmdCtx.Flags)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	env, err := es.EnvCreate(ctx, domain.EnvCreateArgs{
		Name:       mustGetNameArg(cmdCtx.Flags),
		Comment:    commonCreateArgs.Comment,
		CreateTime: commonCreateArgs.CreateTime,
		UpdateTime: commonCreateArgs.UpdateTime,
	})

	if err != nil {
		return fmt.Errorf("could not create env: %w", err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Created env: %s\n", env.Name)

	return nil
}

func EnvDeleteCmd() command.Command {
	return command.New(
		"Delete an environment and associated vars",
		envDeleteRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(confirmFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
	)
}

func envDeleteRun(cmdCtx command.Context) error {

	name := mustGetNameArg(cmdCtx.Flags)

	confirm := mustGetConfirmArg(cmdCtx.Flags)

	if confirm {
		keepGoing, err := askConfirm()
		if err != nil {
			panic(err)
		}
		if !keepGoing {
			return errors.New("unconfirmed change")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer cancel()

	err = es.EnvDelete(ctx, name)

	if err != nil {
		return fmt.Errorf("could not delete env: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "deleted: %s\n", name)
	return nil
}

func EnvListCmd() command.Command {
	return command.New(
		"List environments",
		envListRun,
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
	)
}

func envListRun(cmdCtx command.Context) error {

	timezone := mustGetTimezoneArg(cmdCtx.Flags)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	envs, err := es.EnvList(ctx)
	if err != nil {
		return err
	}

	tableprint.EnvList(cmdCtx.Stdout, envs, tableprint.Timezone(timezone))
	return nil
}

func EnvPrintScriptCmd() command.Command {
	return command.New(
		"Print export script",
		envPrintScriptRun,
		command.ExistingFlag("--name", envNameFlag()),
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
		command.Flag(
			"--shell",
			"Shell to print script for",
			scalar.String(
				scalar.Choices("zsh"),
				scalar.Default("zsh"),
			),
			flag.Required(),
		),
		command.Flag(
			"--type",
			"Type of script",
			scalar.String(
				scalar.Choices("export", "unexport"),
				scalar.Default("export"),
			),
			flag.Required(),
		),
	)
}

func envPrintScriptRun(cmdCtx command.Context) error {
	name := mustGetNameArg(cmdCtx.Flags)
	noEnvNoProblem := cmdCtx.Flags["--no-env-no-problem"].(bool)
	scriptType := cmdCtx.Flags["--type"].(string)
	shell := cmdCtx.Flags["--shell"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	envVars, err := es.EnvVarList(ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrEnvNotFound) && noEnvNoProblem {
			return nil
		}
		return fmt.Errorf("could not list env vars: %s: %w", name, err)
	}

	envRefs, envRefVars, err := es.EnvRefList(ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrEnvNotFound) && noEnvNoProblem {
			return nil
		}
		return fmt.Errorf("could not list env refs: %s: %w", name, err)
	}

	switch shell {
	case "zsh":
		switch scriptType {
		case "export":
			for _, ev := range envVars {
				fmt.Fprintf(cmdCtx.Stdout, "echo 'Adding:' %s;\n", shellescape.Quote(ev.Name))
				fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(ev.Name), shellescape.Quote(ev.Value))
			}

			for i := range len(envRefs) {
				fmt.Fprintf(cmdCtx.Stdout, "echo 'Adding:' %s;\n", shellescape.Quote(envRefs[i].Name))
				fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(envRefs[i].Name), shellescape.Quote(envRefVars[i].Value))
			}
		case "unexport":
			for _, ev := range envVars {
				fmt.Fprintf(cmdCtx.Stdout, "echo 'Removing:' %s;\n", shellescape.Quote(ev.Name))
				fmt.Fprintf(cmdCtx.Stdout, "unset %s;\n", shellescape.Quote(ev.Name))
			}

			for _, er := range envRefs {
				fmt.Fprintf(cmdCtx.Stdout, "echo 'Removing:' %s;\n", shellescape.Quote(er.Name))
				fmt.Fprintf(cmdCtx.Stdout, "unset %s;\n", shellescape.Quote(er.Name))
			}
		default:
			return errors.New("unimplemented --script-type: " + scriptType)

		}
	default:
		return fmt.Errorf("unimplemented shell: %s", shell)
	}

	return nil
}

func EnvShowCmd() command.Command {
	return command.New(
		"Print environment details",
		envShowRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
	)
}

func envShowRun(cmdCtx command.Context) error {

	name := mustGetNameArg(cmdCtx.Flags)
	timezone := mustGetTimezoneArg(cmdCtx.Flags)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	env, err := es.EnvShow(ctx, name)
	if err != nil {
		return fmt.Errorf("could not show env: %s: %w", name, err)
	}

	localvars, err := es.EnvVarList(ctx, name)
	if err != nil {
		return err
	}

	refs, referencedVars, err := es.EnvRefList(ctx, name)
	if err != nil {
		return err
	}

	tableprint.EnvShowRun(cmdCtx.Stdout, *env, localvars, refs, referencedVars, tableprint.Timezone(timezone))
	return nil
}

func EnvUpdateCmd() command.Command {
	return command.New(
		"Update an environment",
		envUpdateRun,
		command.ExistingFlags(commonUpdateFlags()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(confirmFlag()),
	)
}

func envUpdateRun(cmdCtx command.Context) error {

	// common update flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := ptrFromMap[time.Time](cmdCtx.Flags, "--create-time")
	newName := ptrFromMap[string](cmdCtx.Flags, "--new-name")
	updateTime := ptrFromMap[time.Time](cmdCtx.Flags, "--update-time")

	confirm := mustGetConfirmArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)

	if confirm {
		keepGoing, err := askConfirm()
		if err != nil {
			panic(err)
		}
		if !keepGoing {
			return errors.New("unconfirmed change")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	err = es.EnvUpdate(ctx, name, domain.EnvUpdateArgs{
		Comment:    comment,
		CreateTime: createTime,
		NewName:    newName,
		UpdateTime: updateTime,
	})

	if err != nil {
		return fmt.Errorf("could not update env: %w", err)
	}
	finalName := name
	if newName != nil {
		finalName = *newName
	}
	fmt.Fprintln(cmdCtx.Stdout, "updated env:", finalName)
	return nil
}
