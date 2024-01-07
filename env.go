package main

import (
	"fmt"

	"go.bbkane.com/warg/command"
)

func envCreateCmd(cmdCtx command.Context) error {
	name := cmdCtx.Flags["--name"].(string)
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	fmt.Println(name)

	db, err := initSqlite(sqliteDSN)
	if err != nil {
		return fmt.Errorf("could not init db: %w", err)
	}
	fmt.Println(db)

	return nil
}
