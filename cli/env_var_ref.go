package cli

import (
	"fmt"
	"time"

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
		domain.EnvLocalRefCreateArgs{
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
