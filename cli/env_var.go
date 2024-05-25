package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvVarCreateCmd() command.Command {
	return command.New(
		"Create a variable local to the this env",
		envVarCreateRun,
		command.Flag(
			"--value",
			"Value for this local env var",
			scalar.String(),
		),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(commonCreateFlagMap()),
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

func envVarCreateRun(cmdCtx command.Context) error {

	// common create Flags
	commonCreateArgs := mustGetCommonCreateArgs(cmdCtx.Flags)

	envName := mustGetEnvNameArg(cmdCtx.Flags)
	value, exists := cmdCtx.Flags["--value"].(string)
	if !exists {
		fmt.Print("Enter value: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Err() != nil {
			return fmt.Errorf("couldn't read --value: %w", scanner.Err())
		}
		value = scanner.Text()
	}

	name := mustGetNameArg(cmdCtx.Flags)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	_, err = es.EnvVarCreate(
		ctx,
		domain.EnvVarCreateArgs{
			EnvName:    envName,
			Name:       name,
			Comment:    commonCreateArgs.Comment,
			CreateTime: commonCreateArgs.CreateTime,
			UpdateTime: commonCreateArgs.UpdateTime,
			Value:      value,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create env var: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Created env var: %s: %s\n", envName, name)
	return nil
}

func EnvVarDeleteCmd() command.Command {
	return command.New(
		"Delete a variable local to the this env",
		envVarDeleteRun,
		command.ExistingFlags(confirmFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
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

func envVarDeleteRun(cmdCtx command.Context) error {
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

	err = es.EnvVarDelete(ctx, envName, name)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmdCtx.Stdout, "Deleted %s: %s\n", envName, name)
	return nil
}

func EnvVarShowCmd() command.Command {
	return command.New(
		"Show details for a local var",
		envVarShowRun,
		command.ExistingFlags(maskFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.ExistingFlags(formatFlag()),
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

func envVarShowRun(cmdCtx command.Context) error {

	mask := mustGetMaskArg(cmdCtx.Flags)
	envName := mustGetEnvNameArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)
	timezone := mustGetTimezoneArg(cmdCtx.Flags)
	format := cmdCtx.Flags["--format"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), mustGetTimeoutArg(cmdCtx.Flags))
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	envVar, envRefs, err := es.EnvVarShow(ctx, envName, name)
	if err != nil {
		return fmt.Errorf("couldn't find env var: %s: %w", name, err)
	}

	c := tableprint.CommonTablePrintArgs{
		Format: tableprint.Format(format),
		Mask:   mask,
		Tz:     tableprint.Timezone(timezone),
		W:      cmdCtx.Stdout,
	}

	tableprint.EnvLocalVarShowPrint(c, *envVar, envRefs)
	return nil
}
