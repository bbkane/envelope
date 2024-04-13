package cli

import (
	"context"
	"errors"
	"fmt"

	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvRefCreateCmd() command.Command {
	return command.New(
		"Create a reference in this env to a variable in another env",
		envRefCreateRun,
		command.Flag(
			"--name",
			"Ref name",
			scalar.String(),
			flag.Required(),
		),
		command.Flag(
			"--ref-env-name",
			"Environment we're referencing",
			scalar.String(),
			flag.Required(),
		),
		command.Flag(
			"--ref-var-name",
			"Variable we're referencing",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlags(commonCreateFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
	)
}

func envRefCreateRun(cmdCtx command.Context) error {
	// common create Flags
	commonCreateArgs := mustGetCommonCreateArgs(cmdCtx.Flags)

	name := mustGetNameArg(cmdCtx.Flags)
	refEnvName := cmdCtx.Flags["--ref-env-name"].(string)
	refVarName := cmdCtx.Flags["--ref-var-name"].(string)

	envName := mustGetEnvNameArg(cmdCtx.Flags)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	_, err = es.EnvRefCreate(
		ctx,
		domain.EnvRefCreateArgs{
			EnvName:    envName,
			Name:       name,
			Comment:    commonCreateArgs.Comment,
			CreateTime: commonCreateArgs.CreateTime,
			UpdateTime: commonCreateArgs.UpdateTime,
			RefEnvName: refEnvName,
			RefVarName: refVarName,
		},
	)

	if err != nil {
		return fmt.Errorf("couldn't create env ref: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Created env ref: %s: %s\n", envName, name)
	return nil
}

func EnvRefDeleteCmd() command.Command {
	return command.New(
		"Delete a reference to a variablea",
		envRefDeleteRun,
		command.ExistingFlags(confirmFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.Flag(
			"--name",
			"Ref name",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
	)
}

func envRefDeleteRun(cmdCtx command.Context) error {
	envName := mustGetEnvNameArg(cmdCtx.Flags)

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

	err = es.EnvRefDelete(ctx, envName, name)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmdCtx.Stdout, "Deleted %s: %s\n", envName, name)
	return nil
}

func EnvRefShowCmd() command.Command {
	return command.New(
		"Show details for a reference",
		envRefShowRun,
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.Flag(
			"--name",
			"Env ref name",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
	)
}

func envRefShowRun(cmdCtx command.Context) error {
	envName := mustGetEnvNameArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)
	timezone := mustGetTimezoneArg(cmdCtx.Flags)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	envRef, envVar, err := es.EnvRefShow(ctx, envName, name)
	if err != nil {
		return fmt.Errorf("couldn't find env var: %s: %w", name, err)
	}

	tableprint.EnvRefShowPrint(cmdCtx.Stdout, *envRef, *envVar, tableprint.Timezone(timezone))
	return nil
}
