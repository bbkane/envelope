package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/alessio/shellescape"
	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/tableprint"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvCreateCmd() command.Command {
	return command.New(
		"Create an environment",
		envCreateRun,
		command.ExistingFlags(commonCreateFlag()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
	)
}

func envCreateRun(cmdCtx command.Context) error {

	// common create Flags
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	name := cmdCtx.Flags["--name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	env, err := iesr.EnvService.EnvCreate(iesr.Ctx, domain.EnvCreateArgs{
		Name:       name,
		Comment:    comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
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
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
	)
}

func envDeleteRun(cmdCtx command.Context) error {

	name := cmdCtx.Flags["--name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	err = iesr.EnvService.EnvDelete(iesr.Ctx, name)

	if err != nil {
		return fmt.Errorf("could not delete env: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "deleted: %s\n", name)
	return nil
}

func EnvPrintScriptCmd() command.Command {
	return command.New(
		"Print export script",
		envPrintScriptRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
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
	name := cmdCtx.Flags["--name"].(string)
	scriptType := cmdCtx.Flags["--type"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	envVars, err := iesr.EnvService.EnvVarLocalList(iesr.Ctx, name)
	if err != nil {
		return fmt.Errorf("could not list env vars: %s: %w", name, err)
	}
	if scriptType == "export" {
		for _, ev := range envVars {
			fmt.Fprintf(cmdCtx.Stdout, "echo 'Adding:' %s;\n", shellescape.Quote(ev.Name))
			fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(ev.Name), shellescape.Quote(ev.Value))
		}
	} else {
		return errors.New("Unimplemented --script-type: " + scriptType)
	}

	return nil
}

func EnvShowCmd() command.Command {
	return command.New(
		"Print environment details",
		envShowRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
		command.ExistingFlags(timeZoneFlagMap()),
	)
}

func envShowRun(cmdCtx command.Context) error {

	name := cmdCtx.Flags["--name"].(string)
	timezone := cmdCtx.Flags["--timezone"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	env, err := iesr.EnvService.EnvShow(iesr.Ctx, name)
	if err != nil {
		return fmt.Errorf("could not show env: %s: %w", name, err)
	}

	tableprint.EnvTable(cmdCtx.Stdout, *env, tableprint.Timezone(timezone))
	return nil
}

func EnvUpdateCmd() command.Command {
	return command.New(
		"Update an environment",
		envUpdateRun,
		command.ExistingFlags(commonUpdateFlags()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
	)
}

func envUpdateRun(cmdCtx command.Context) error {

	// common update flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := ptrFromMap[time.Time](cmdCtx.Flags, "--create-time")
	newName := ptrFromMap[string](cmdCtx.Flags, "--new-name")
	updateTime := ptrFromMap[time.Time](cmdCtx.Flags, "--update-time")

	name := cmdCtx.Flags["--name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	err = iesr.EnvService.EnvUpdate(iesr.Ctx, name, domain.EnvUpdateArgs{
		Comment:    comment,
		CreateTime: createTime,
		NewName:    newName,
		UpdateTime: updateTime,
	})

	if err != nil {
		return fmt.Errorf("could not update env: %w", err)
	}
	return nil
}
