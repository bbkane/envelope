package cli

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/keyring"
	"go.bbkane.com/namedenv/sqlite"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvLocalVarCreateCmd() command.Command {
	return command.New(
		"Create a variable local to the this env",
		envVarCreateLocalRun,
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

func envVarCreateLocalRun(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common create Flags
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	envName := cmdCtx.Flags["--env-name"].(string)
	value := cmdCtx.Flags["--value"].(string)

	name := cmdCtx.Flags["--name"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	envVar, err := envService.EnvVarLocalCreate(
		ctx,
		domain.EnvVarLocalCreateArgs{
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

	fmt.Fprintf(cmdCtx.Stdout, "Created env var: %#v\n", envVar)
	return nil

}
