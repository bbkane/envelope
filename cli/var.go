package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/models"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func VarCreateCmd() command.Command {
	return command.New(
		"Create a variable local to the this env",
		withEnvService(varCreateRun),
		command.ExistingFlag(
			"--env-name",
			envNameFlag(),
		),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(commonCreateFlagMap()),
		command.Flag(
			"--name",
			"Existing env var name",
			scalar.String(),
			flag.Required(),
		),
		command.Flag(
			"--value",
			"Value for this local env var",
			scalar.String(),
		),
	)
}

func varCreateRun(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {

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

	_, err := es.VarCreate(
		ctx,
		models.VarCreateArgs{
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

func VarDeleteCmd() command.Command {
	return command.New(
		"Delete a variable local to the this env",
		withConfirm(withEnvService(varDeleteRun)),
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

func varDeleteRun(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	envName := mustGetEnvNameArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)

	err := es.VarDelete(ctx, envName, name)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmdCtx.Stdout, "Deleted %s: %s\n", envName, name)
	return nil
}

func VarShowCmd() command.Command {
	return command.New(
		"Show details for a local var",
		withEnvService(varShowRun),
		command.ExistingFlags(maskFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.ExistingFlags(formatFlag()),
		command.ExistingFlags(widthFlag()),
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

func varShowRun(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {

	mask := mustGetMaskArg(cmdCtx.Flags)
	envName := mustGetEnvNameArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)
	timezone := mustGetTimezoneArg(cmdCtx.Flags)
	format := cmdCtx.Flags["--format"].(string)
	width := mustGetWidthArg(cmdCtx.Flags)

	envVar, envRefs, err := es.VarShow(ctx, envName, name)
	if err != nil {
		return fmt.Errorf("couldn't find env var: %s: %w", name, err)
	}

	c := tableprint.CommonTablePrintArgs{
		Format:          tableprint.Format(format),
		Mask:            mask,
		Tz:              tableprint.Timezone(timezone),
		W:               cmdCtx.Stdout,
		DesiredMaxWidth: width,
	}

	tableprint.VarShowPrint(c, *envVar, envRefs)
	return nil
}

func VarUpdateCmd() command.Command {
	return command.New(
		"Update and env var",
		withConfirm(withEnvService(varUpdateRun)),
		command.ExistingFlag("--env-name", envNameFlag()),
		command.ExistingFlags(commonUpdateFlags()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(confirmFlag()),
		command.Flag(
			"--name",
			"Env var name",
			scalar.String(),
			flag.Required(),
		),
		command.Flag(
			"--new-env-name",
			"New env name",
			scalar.String(),
		),
		command.Flag(
			"--value",
			"New value for this env var",
			scalar.String(),
		),
	)
}

func varUpdateRun(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	// common update flags
	commonUpdateArgs := getCommonUpdateArgs(cmdCtx.Flags)

	envName := mustGetEnvNameArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)
	newEnvName := ptrFromMap[string](cmdCtx.Flags, "--new-env-name")
	value := ptrFromMap[string](cmdCtx.Flags, "--value")

	err := es.VarUpdate(ctx, envName, name, models.VarUpdateArgs{
		Comment:    commonUpdateArgs.Comment,
		CreateTime: commonUpdateArgs.CreateTime,
		EnvName:    newEnvName,
		Name:       commonUpdateArgs.NewName,
		UpdateTime: commonUpdateArgs.UpdateTime,
		Value:      value,
	})

	if err != nil {
		return fmt.Errorf("could not update env var: %w", err)
	}
	finalName := name
	if commonUpdateArgs.NewName != nil {
		finalName = *commonUpdateArgs.NewName
	}
	finalEnvName := envName
	if newEnvName != nil {
		finalEnvName = *newEnvName
	}
	fmt.Fprintf(cmdCtx.Stdout, "updated env var:  %s: %s\n", finalEnvName, finalName)
	return nil
}
