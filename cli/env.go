package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/domain"

	"go.bbkane.com/warg/command"
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
		command.ExistingFlags(widthFlag()),
	)
}

func envListRun(cmdCtx command.Context) error {

	timezone := mustGetTimezoneArg(cmdCtx.Flags)
	width := mustGetWidthArg(cmdCtx.Flags)

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

	c := tableprint.CommonTablePrintArgs{
		Format:          tableprint.Format_Table,
		Mask:            false,
		Tz:              tableprint.Timezone(timezone),
		W:               cmdCtx.Stdout,
		DesiredMaxWidth: width,
	}

	tableprint.EnvList(c, envs)
	return nil
}

func EnvShowCmd() command.Command {
	return command.New(
		"Print environment details",
		envShowRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(maskFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.ExistingFlags(widthFlag()),
	)
}

func envShowRun(cmdCtx command.Context) error {

	mask := mustGetMaskArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)
	timezone := mustGetTimezoneArg(cmdCtx.Flags)
	width := mustGetWidthArg(cmdCtx.Flags)

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

	c := tableprint.CommonTablePrintArgs{
		// NOTE: since the only two options are value-only and table,
		// and value-only doesn't make sense here, hardcode Format_Table
		Format:          tableprint.Format_Table,
		Mask:            mask,
		Tz:              tableprint.Timezone(timezone),
		W:               cmdCtx.Stdout,
		DesiredMaxWidth: width,
	}
	tableprint.EnvShowRun(c, *env, localvars, refs, referencedVars)
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
	timeout := mustGetTimeoutArg(cmdCtx.Flags)

	if confirm {
		keepGoing, err := askConfirm()
		if err != nil {
			panic(err)
		}
		if !keepGoing {
			return errors.New("unconfirmed change")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	es, err := initEnvService(ctx, cmdCtx.Flags)
	if err != nil {
		return err
	}

	err = es.EnvUpdate(ctx, name, domain.EnvUpdateArgs{
		Comment:    comment,
		CreateTime: createTime,
		Name:       newName,
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
