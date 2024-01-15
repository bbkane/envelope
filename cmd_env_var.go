package main

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite"
	"go.bbkane.com/warg/command"
)

func envVarCreateLocalCmd(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common create Flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	envName := cmdCtx.Flags["--env-name"].(string)
	value := cmdCtx.Flags["--value"].(string)

	name := cmdCtx.Flags["--name"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	envVar, err := envService.CreateLocalEnvVar(
		ctx,
		domain.CreateLocalEnvVarArgs{
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
