package main

import (
	"context"
	"fmt"
	"time"

	"go.bbkane.com/namedenv/domain"
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
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)
	name := cmdCtx.Flags["--name"].(string)
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

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

	return nil
}

func envUpdateCmd(cmdCtx command.Context) error {
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := ptrFromMap[time.Time](cmdCtx.Flags, "--create-time")
	name := cmdCtx.Flags["--name"].(string)
	newName := ptrFromMap[string](cmdCtx.Flags, "--new-name")
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)
	updateTime := ptrFromMap[time.Time](cmdCtx.Flags, "--update-time")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	id, err := envService.FindEnvByName(ctx, name)
	if err != nil {
		return err
	}

	err = envService.UpdateEnv(ctx, id, domain.UpdateEnvArgs{
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
