package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/alessio/shellescape"
	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/keyring"
	"go.bbkane.com/namedenv/sqlite"
	"go.bbkane.com/namedenv/tableprint"

	"go.bbkane.com/warg/command"
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
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common create Flags
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	name := cmdCtx.Flags["--name"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	env, err := envService.EnvCreate(ctx, domain.EnvCreateArgs{
		Name:       name,
		Comment:    comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
	})

	if err != nil {
		return fmt.Errorf("could not create env: %w", err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Created env: %#v\n", env)

	return nil
}

func EnvUpdateCmd() command.Command {
	return command.New(
		"Update an environment",
		EnvUpdateRun,
		command.ExistingFlags(commonUpdateFlags()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
	)
}

func EnvUpdateRun(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common update flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := ptrFromMap[time.Time](cmdCtx.Flags, "--create-time")
	newName := ptrFromMap[string](cmdCtx.Flags, "--new-name")
	updateTime := ptrFromMap[time.Time](cmdCtx.Flags, "--update-time")

	name := cmdCtx.Flags["--name"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	err = envService.EnvUpdate(ctx, name, domain.EnvUpdateArgs{
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

func EnvPrintScriptExportCmd() command.Command {
	return command.New(
		"Print export script",
		EnvPrintScriptExportRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
	)
}

func EnvPrintScriptExportRun(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	name := cmdCtx.Flags["--name"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	envVars, err := envService.EnvVarLocalList(ctx, name)
	if err != nil {
		return fmt.Errorf("could not list env vars: %s: %w", name, err)
	}

	for _, ev := range envVars {
		fmt.Fprintf(cmdCtx.Stdout, "echo 'Adding:' %s;\n", shellescape.Quote(ev.Name))
		fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(ev.Name), shellescape.Quote(ev.Value))
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
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	name := cmdCtx.Flags["--name"].(string)
	timezone := cmdCtx.Flags["--timezone"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	env, err := envService.EnvShow(ctx, name)
	if err != nil {
		return fmt.Errorf("could not show env: %s: %w", name, err)
	}

	tableprint.EnvTable(cmdCtx.Stdout, *env, tableprint.Timezone(timezone))
	return nil
}
