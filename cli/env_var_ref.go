package cli

import (
	"errors"
	"fmt"
	"time"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/tableprint"
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
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	name := cmdCtx.Flags["--name"].(string)
	refEnvName := cmdCtx.Flags["--ref-env-name"].(string)
	refVarName := cmdCtx.Flags["--ref-var-name"].(string)

	envName := cmdCtx.Flags["--env-name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	_, err = iesr.EnvService.EnvRefCreate(
		iesr.Ctx,
		domain.EnvRefCreateArgs{
			EnvName:    envName,
			Name:       name,
			Comment:    comment,
			CreateTime: createTime,
			UpdateTime: updateTime,
			RefEnvName: refEnvName,
			RefVarName: refVarName,
		},
	)

	if err != nil {
		return fmt.Errorf("couldn't creaate env ref: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Creeated env ref: %s: %s\n", envName, name)
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
	envName := cmdCtx.Flags["--env-name"].(string)

	confirm := cmdCtx.Flags["--confirm"].(bool)
	name := cmdCtx.Flags["--name"].(string)

	if confirm {
		keepGoing, err := askConfirm()
		if err != nil {
			panic(err)
		}
		if !keepGoing {
			return errors.New("unconfirmed change")
		}
	}

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	err = iesr.EnvService.EnvRefDelete(iesr.Ctx, envName, name)
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
	envName := cmdCtx.Flags["--env-name"].(string)
	name := cmdCtx.Flags["--name"].(string)
	timezone := cmdCtx.Flags["--timezone"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	envRef, envVar, err := iesr.EnvService.EnvRefShow(iesr.Ctx, envName, name)
	if err != nil {
		return fmt.Errorf("couldn't find env var: %s: %w", name, err)
	}

	tableprint.EnvRefShowPrint(cmdCtx.Stdout, *envRef, *envVar, tableprint.Timezone(timezone))
	return nil
}
