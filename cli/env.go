package cli

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/models"

	"go.bbkane.com/warg/command"
)

func EnvCreateCmd() command.Command {
	return command.New(
		"Create an environment",
		withEnvService(envCreate),
		command.ExistingFlags(commonCreateFlagMap()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
	)
}

func envCreate(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	commonCreateArgs := mustGetCommonCreateArgs(cmdCtx.Flags)

	env, err := es.EnvCreate(ctx, models.EnvCreateArgs{
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
		withConfirm(withEnvService(envDelete)),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(confirmFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
	)
}

func envDelete(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	name := mustGetNameArg(cmdCtx.Flags)

	err := es.EnvDelete(ctx, name)
	if err != nil {
		return fmt.Errorf("could not delete env: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "deleted: %s\n", name)
	return nil
}

func EnvListCmd() command.Command {
	return command.New(
		"List environments",
		withEnvService(envList),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.ExistingFlags(widthFlag()),
	)
}

func envList(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	envs, err := es.EnvList(ctx)
	if err != nil {
		return err
	}

	c := tableprint.CommonTablePrintArgs{
		Format:          tableprint.Format_Table,
		Mask:            false,
		Tz:              tableprint.Timezone(mustGetTimezoneArg(cmdCtx.Flags)),
		W:               cmdCtx.Stdout,
		DesiredMaxWidth: mustGetWidthArg(cmdCtx.Flags),
	}

	tableprint.EnvList(c, envs)
	return nil
}

func EnvShowCmd() command.Command {
	return command.New(
		"Print environment details",
		withEnvService(envShow),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(maskFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
		command.ExistingFlags(widthFlag()),
	)
}

func envShow(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	mask := mustGetMaskArg(cmdCtx.Flags)
	name := mustGetNameArg(cmdCtx.Flags)
	timezone := mustGetTimezoneArg(cmdCtx.Flags)
	width := mustGetWidthArg(cmdCtx.Flags)

	env, err := es.EnvShow(ctx, name)
	if err != nil {
		return fmt.Errorf("could not show env: %s: %w", name, err)
	}

	localvars, err := es.VarList(ctx, name)
	if err != nil {
		return err
	}

	refs, referencedVars, err := es.VarRefList(ctx, name)
	if err != nil {
		return err
	}

	c := tableprint.CommonTablePrintArgs{
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
		withConfirm(withEnvService(envUpdate)),
		command.ExistingFlags(commonUpdateFlags()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(confirmFlag()),
	)
}

func envUpdate(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
	// common update flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := ptrFromMap[time.Time](cmdCtx.Flags, "--create-time")
	newName := ptrFromMap[string](cmdCtx.Flags, "--new-name")
	updateTime := ptrFromMap[time.Time](cmdCtx.Flags, "--update-time")

	name := mustGetNameArg(cmdCtx.Flags)

	err := es.EnvUpdate(ctx, name, models.EnvUpdateArgs{
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
