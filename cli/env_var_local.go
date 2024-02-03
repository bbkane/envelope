package cli

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/tableprint"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvLocalVarCreateCmd() command.Command {
	return command.New(
		"Create a variable local to the this env",
		envLocalVarCreateRun,
		command.Flag(
			"--value",
			"Value for this local env var",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
		command.ExistingFlags(commonCreateFlag()),
		command.Flag(
			"--name",
			"Env var name",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
	)
}

func envLocalVarCreateRun(cmdCtx command.Context) error {

	// common create Flags
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	envName := cmdCtx.Flags["--env-name"].(string)
	value := cmdCtx.Flags["--value"].(string)

	name := cmdCtx.Flags["--name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	_, err = iesr.EnvService.EnvLocalVarCreate(
		iesr.Ctx,
		domain.EnvLocalVarCreateArgs{
			EnvName:    envName,
			Name:       name,
			Comment:    comment,
			CreateTime: createTime,
			UpdateTime: updateTime,
			Value:      value,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create env var: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Created env var: %s: %s\n", envName, name)
	return nil
}

func EnvLocalVarDeleteCmd() command.Command {
	return command.New(
		"Delete a variable local to the this env",
		envLocalVarDeleteRun,
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
		command.Flag(
			"--name",
			"Env var name",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
	)
}

func envLocalVarDeleteRun(cmdCtx command.Context) error {
	envName := cmdCtx.Flags["--env-name"].(string)

	name := cmdCtx.Flags["--name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	err = iesr.EnvService.EnvLocalVarDelete(context.TODO(), envName, name)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmdCtx.Stdout, "Deleted %s: %s\n", envName, name)
	return nil
}

func EnvLocalVarShowCmd() command.Command {
	return command.New(
		"Show details for a localvar",
		envLocalVarShowRun,
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.Flag(
			"--name",
			"Env var name",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
	)
}

func envLocalVarShowRun(cmdCtx command.Context) error {

	envName := cmdCtx.Flags["--env-name"].(string)
	name := cmdCtx.Flags["--name"].(string)
	timezone := cmdCtx.Flags["--timezone"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	envVar, err := iesr.EnvService.EnvLocalVarShow(iesr.Ctx, envName, name)
	if err != nil {
		return fmt.Errorf("couldn't find env var: %s: %w", name, err)
	}

	tableprint.EnvLocalVarShowPrint(cmdCtx.Stdout, *envVar, tableprint.Timezone(timezone))
	return nil
}