package main

import (
	"context"
	"fmt"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite"

	"go.bbkane.com/warg/command"
)

func envCreateCmd(cmdCtx command.Context) error {
	createTimeStr := cmdCtx.Flags["--create-time"].(string)
	updateTimeStr := cmdCtx.Flags["--update-time"].(string)
	name := cmdCtx.Flags["--name"].(string)
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)

	var comment *string
	commentFlag, commentExists := cmdCtx.Flags["--comment"]
	if commentExists {
		tmp := commentFlag.(string)
		comment = &tmp
	}

	createTime, err := domain.StringToTime(createTimeStr)
	if err != nil {
		return fmt.Errorf("could not parse --create-time: %w", err)
	}
	updateTime, err := domain.StringToTime(updateTimeStr)
	if err != nil {
		return fmt.Errorf("could not parse --update-time: %w", err)
	}

	ctx := context.Background() // TODO: fix
	envService, err := sqlite.NewEnvService(ctx, sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	envID, err := envService.CreateEnv(ctx, domain.CreateEnvArgs{
		Name:       name,
		Comment:    comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
	})

	if err != nil {
		return fmt.Errorf("could not create env: %w", err)
	}

	fmt.Printf("Created env with ID: %#v\n", envID)

	// TODO: make this it's own command...
	newName := name + name
	err = envService.UpdateEnv(ctx, domain.UpdateEnvArgs{
		Name: &newName,
		ID:   envID,
	})

	if err != nil {
		return fmt.Errorf("TODO: rm me but here's the error: %w", err)
	}

	return nil
}
