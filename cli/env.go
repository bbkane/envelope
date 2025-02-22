package cli

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/envelope/cli/tableprint"
	"go.bbkane.com/envelope/models"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvCreateCmd() command.Command {
	return command.New(
		"Create an environment",
		withEnvService(envCreate),
		command.FlagMap(commonCreateFlagMap()),
		command.Flag("--name", envNameFlag()),
		command.FlagMap(timeoutFlagMap()),
		command.FlagMap(sqliteDSNFlagMap()),
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
		command.Flag("--name", envNameFlag()),
		command.FlagMap(confirmFlag()),
		command.FlagMap(timeoutFlagMap()),
		command.FlagMap(sqliteDSNFlagMap()),
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
		command.FlagMap(timeoutFlagMap()),
		command.FlagMap(sqliteDSNFlagMap()),
		command.FlagMap(timeZoneFlagMap()),
		command.FlagMap(widthFlag()),
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
		command.Flag("--name", envNameFlag()),
		command.FlagMap(maskFlag()),
		command.FlagMap(timeoutFlagMap()),
		command.FlagMap(sqliteDSNFlagMap()),
		command.FlagMap(timeZoneFlagMap()),
		command.FlagMap(widthFlag()),
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

func EnvUpdateCmd2() command.Command {
	// TODO: rework this to use the new PointerTo system. Make the helper functions less awkward
	var envName string
	var updateArgs models.EnvUpdateArgs
	return command.New(
		"Update an environment",
		withConfirm(withEnvService(func(ctx context.Context, es models.EnvService, cmdCtx command.Context) error {
			return nil
		})),
		command.FlagMap(flag.FlagMap{
			"--comment": flag.New(
				"Comment",
				scalar.String(
					scalar.PointerTo(updateArgs.Comment),
				),
			),
			"--create-time": flag.New(
				"Create time",
				scalar.New(
					datetime(),
					scalar.PointerTo(updateArgs.CreateTime),
				),
			),
			"--new-name": flag.New(
				"New name",
				scalar.String(
					scalar.PointerTo(updateArgs.Name),
				),
			),
			"--update-time": flag.New(
				"Update time",
				scalar.New(
					datetime(),
					scalar.Default(time.Now()),
					scalar.PointerTo(updateArgs.UpdateTime),
				),
				flag.UnsetSentinel("UNSET"),
			),
		}),
		// NOTE: this is to find the env
		command.Flag("--name", envNameFlag2(&envName)),
		command.FlagMap(timeoutFlagMap()),
		command.FlagMap(sqliteDSNFlagMap()),
		command.FlagMap(confirmFlag()),
	)
}

func EnvUpdateCmd() command.Command {
	return command.New(
		"Update an environment",
		withConfirm(withEnvService(envUpdate)),
		command.FlagMap(commonUpdateFlags()),
		command.Flag("--name", envNameFlag()),
		command.FlagMap(timeoutFlagMap()),
		command.FlagMap(sqliteDSNFlagMap()),
		command.FlagMap(confirmFlag()),
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
