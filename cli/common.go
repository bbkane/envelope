package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/keyring"
	"go.bbkane.com/envelope/sqlite"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"

	"go.bbkane.com/warg/value/contained"
	"go.bbkane.com/warg/value/scalar"
)

func askConfirm() (bool, error) {
	fmt.Print("Type 'yes' to continue: ")
	reader := bufio.NewReader(os.Stdin)
	confirmation, err := reader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("confirmation ReadString error: %w", err)
		return false, err
	}
	confirmation = strings.TrimSpace(confirmation)
	if confirmation != "yes" {
		return false, nil
	}
	return true, nil
}

func emptyOrNil[T any](iFace interface{}) (T, error) {
	under, ok := iFace.(T)
	if !ok {
		return under, contained.ErrIncompatibleInterface
	}
	return under, nil
}

func datetime() contained.TypeInfo[time.Time] {
	return contained.TypeInfo[time.Time]{
		Description: "datetime in RFC3339 format",
		FromIFace:   emptyOrNil[time.Time],
		FromInstance: func(t time.Time) (time.Time, error) {
			return t, nil
		},
		FromString: func(s string) (time.Time, error) {
			return time.Parse(time.RFC3339, s)
		},
		Empty: func() time.Time {
			return time.Time{}
		},
	}
}

func confirmFlag() flag.FlagMap {
	return flag.FlagMap{
		"--confirm": flag.New(
			"Ask for confirmation before running",
			scalar.Bool(
				scalar.Default(true),
			),
			flag.Required(),
		),
	}
}

func envNameFlag() flag.Flag {

	cwd, err := os.Getwd()
	if err != nil {
		// I don't know when this could happen?
		panic(err)
	}

	envNameFlag := flag.New(
		"Environment name",
		scalar.String(
			scalar.Default(cwd),
		),
		flag.Required(),
	)
	return envNameFlag
}

func sqliteDSNFlagMap() flag.FlagMap {
	dbPath, err := homedir.Expand("~/.config/envelope.db")
	if err != nil {
		panic(err)
	}

	sqliteDSN := flag.FlagMap{
		"--db-path": flag.New(
			"Sqlite DSN. Usually the file name",
			scalar.String(
				scalar.Default(dbPath),
			),
			flag.Required(),
			flag.EnvVars("ENVELOPE_DB_PATH"),
		),
	}
	return sqliteDSN
}

func commonCreateFlagMap() flag.FlagMap {
	commonCreateFlags := flag.FlagMap{
		"--comment": flag.New(
			"Comment",
			scalar.String(
				scalar.Default(""),
			),
			flag.Required(),
		),
		"--create-time": flag.New(
			"Create time",
			scalar.New(
				datetime(),
				scalar.Default(time.Now()),
			),
			flag.Required(),
		),
		"--update-time": flag.New(
			"Update time",
			scalar.New(
				datetime(),
				scalar.Default(time.Now()),
			),
			flag.Required(),
		),
	}
	return commonCreateFlags
}

func commonUpdateFlags() flag.FlagMap {

	cwd, err := os.Getwd()
	if err != nil {
		// I don't know when this could happen?
		panic(err)
	}

	commonUpdateFlags := flag.FlagMap{
		"--comment": flag.New(
			"Comment",
			scalar.String(),
		),
		"--create-time": flag.New(
			"Create time",
			scalar.New(
				datetime(),
			),
		),
		"--new-name": flag.New(
			"New name",
			scalar.String(
				scalar.Default(cwd),
			),
		),
		"--update-time": flag.New(
			"Update time",
			scalar.New(
				datetime(),
				scalar.Default(time.Now()),
			),
			flag.UnsetSentinel("UNSET"),
		),
	}
	return commonUpdateFlags
}

func timeoutFlagMap() flag.FlagMap {
	timeoutFlag := flag.FlagMap{
		"--timeout": flag.New(
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(10*time.Minute),
			),
			flag.Required(),
		),
	}
	return timeoutFlag

}

func timeZoneFlagMap() flag.FlagMap {
	return flag.FlagMap{
		"--timezone": flag.New(
			"Timezone to display dates",
			scalar.String(
				scalar.Default("local"),
				scalar.Choices("local", "utc"),
			),
			flag.Required(),
		),
	}
}

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

type initEnvServiceRet struct {
	Cancel     context.CancelFunc
	Ctx        context.Context
	EnvService domain.EnvService
}

// initEnvService reads --db-path and --timeout to create a
// EnvService. It means I don't have to type these EVERY time
func initEnvService(passedFlags command.PassedFlags) (*initEnvServiceRet, error) {
	// common flags
	sqliteDSN := passedFlags["--db-path"].(string)
	timeout := passedFlags["--timeout"].(time.Duration)

	//nolint:govet // we don't need the cancel if we err out
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	keyring := keyring.NewOSKeyring(sqliteDSN)

	envService, err := sqlite.NewEnvService(ctx, sqliteDSN, keyring)
	if err != nil {
		//nolint:govet // we don't need the cancel if we err out
		return nil, fmt.Errorf("could not create env service: %w", err)
	}
	return &initEnvServiceRet{
		EnvService: envService,
		Cancel:     cancel,
		Ctx:        ctx,
	}, nil
}
