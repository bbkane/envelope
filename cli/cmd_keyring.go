package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.bbkane.com/namedenv/domain"
	"go.bbkane.com/namedenv/keyring"
	"go.bbkane.com/namedenv/sqlite"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func promptKeyringValue() (string, error) {
	fmt.Print("Enter value to store in keyring: ")
	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	val = strings.TrimSpace(val)
	return val, nil
}

func KeyringCreateCmd() command.Command {
	return command.New(
		"Create a keyring entry. Prompts for value instead of using a flag",
		keyringCreateRun,
		command.ExistingFlags(commonCreateFlag()),
		command.Flag(
			"--name",
			"Keyring entry name",
			scalar.String(),
			flag.Required(),
		),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlag()),
	)
}

func keyringCreateRun(cmdCtx command.Context) error {
	// common flags
	sqliteDSN := cmdCtx.Flags["--sqlite-dsn"].(string)
	timeout := cmdCtx.Flags["--timeout"].(time.Duration)

	// common create Flags
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	name := cmdCtx.Flags["--name"].(string)

	val, err := promptKeyringValue()
	if err != nil {
		return fmt.Errorf("promptKeyringValue err: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		return fmt.Errorf("could not create env service: %w", err)
	}

	entry, err := envService.KeyringEntryCreate(ctx, domain.KeyringEntryCreateArgs{
		Name:       name,
		Comment:    comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
		Value:      val,
	})

	if err != nil {
		return fmt.Errorf("could not create keyring: %w", err)
	}

	// TODO: don't print the value?
	fmt.Fprintf(cmdCtx.Stdout, "Created keyring: %#v\n", entry)

	return nil
}
