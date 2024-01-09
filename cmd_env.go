package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/sqlite"
	"go.bbkane.com/namedenv/sqlite/connect"
	"go.bbkane.com/namedenv/sqlite/sqlcgen"

	"go.bbkane.com/warg/command"
)

func envCreateCmd(cmdCtx command.Context) error {
	// createTime := cmdCtx.Flags["--create-time"].(string)
	// updateTime := cmdCtx.Flags["--update-time"].(string)
	name := cmdCtx.Flags["--name"].(string)
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)

	var comment *string
	commentFlag, commentExists := cmdCtx.Flags["--comment"]
	if commentExists {
		tmp := commentFlag.(string)
		comment = &tmp
	}

	ctx := context.Background() // TODO: fix
	envService, err := sqlite.NewEnvService(ctx, sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	envID, err := envService.CreateEnv(ctx, domain.CreateEnvArgs{
		Name:       name,
		Comment:    comment,
		CreateTime: time.Now(), // TODO: read from flag
		UpdateTime: time.Now(), // TODO: read from flag
	})

	if err != nil {
		return fmt.Errorf("could not create env: %w", err)
	}

	fmt.Printf("Created env with ID: %#v\n", envID)

	return nil
}

func envCreateCmdOld(cmdCtx command.Context) error {
	createTime := cmdCtx.Flags["--create-time"].(string)
	updateTime := cmdCtx.Flags["--update-time"].(string)
	name := cmdCtx.Flags["--name"].(string)
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)

	var comment string
	var commentExists bool
	commentFlag, commentExists := cmdCtx.Flags["--comment"]
	if commentExists {
		comment = commentFlag.(string)
	}

	// TODO: make this use a context
	db, err := connect.Connect(sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not init db: %w", err)
	}

	ctx := context.Background() // TODO: use a timeout!

	queries := sqlcgen.New(db)

	createdEnvID, err := queries.CreateEnv(ctx, sqlcgen.CreateEnvParams{
		Name: name,
		Comment: sql.NullString{
			String: comment,
			Valid:  commentExists,
		},
		CreateTime: createTime,
		UpdateTime: updateTime,
	})
	if err != nil {
		return fmt.Errorf("could not create env: %w", err)
	}

	fmt.Printf("Created env with ID: %#v\n", createdEnvID)

	return nil
}
