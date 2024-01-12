package main

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite"
	"go.bbkane.com/warg/command"
)

func envVarCreateCmd(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common create Flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	envName := cmdCtx.Flags["--env-name"].(string)
	localValue := ptrFromMap[string](cmdCtx.Flags, "--local-value")
	name := cmdCtx.Flags["--name"].(string)
	typeFlg := cmdCtx.Flags["--type"].(string)

	if localValue == nil {
		panic("TODO, need to confirm type and values work")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	envVar, err := envService.CreateEnvVar(
		ctx,
		domain.CreateEnvVarArgs{
			EnvName:    envName,
			Name:       name,
			Comment:    comment,
			CreateTime: createTime,
			UpdateTime: updateTime,
			Type:       domain.EnvVarType(typeFlg),
			LocalValue: localValue,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create env var: %s: %w", name, err)
	}

	fmt.Printf("Created env var: %#v\n", envVar)
	return nil

}
