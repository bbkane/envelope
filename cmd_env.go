package main

import (
	"context"
	"database/sql"
	"fmt"

	"go.bbkane.com/namedenv/envsqlc"
	"go.bbkane.com/namedenv/sqlite/connect"

	"go.bbkane.com/warg/command"
)

func envCreateCmd(cmdCtx command.Context) error {
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

	queries := envsqlc.New(db)

	createdEnvID, err := queries.CreateEnv(ctx, envsqlc.CreateEnvParams{
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
