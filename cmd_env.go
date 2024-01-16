package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alessio/shellescape"
	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/keyring"
	"go.bbkane.com/namedenv/sqlite"

	"go.bbkane.com/warg/command"
)

// ptrFromMap returns &val if key is in the map, otherwise nil
// useful for converting from the cmdCtx.Flags to the types domain needs
func ptrFromMap[T any](m map[string]any, key string) *T {
	val, exists := m[key]
	if exists {
		ret := val.(T)
		return &ret
	}
	return nil
}

func envCreateCmd(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common create Flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
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

	env, err := envService.EnvCreate(ctx, domain.CreateEnvArgs{
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

func envUpdateCmd(cmdCtx command.Context) error {
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

	err = envService.EnvUpdate(ctx, name, domain.UpdateEnvArgs{
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

func envPrintScriptExportCmd(cmdCtx command.Context) error {
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
