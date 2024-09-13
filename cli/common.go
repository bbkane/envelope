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
	"go.bbkane.com/envelope/sqlite"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"golang.org/x/term"

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

// datetime is a type for the CLI so I can pass strings in and parse them to dates
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

func maskFlag() flag.FlagMap {
	return flag.FlagMap{
		"--mask": flag.New(
			"Mask values when printing",
			scalar.Bool(
				scalar.Default(true),
			),
			flag.Required(),
		),
	}
}

func formatFlag() flag.FlagMap {
	return flag.FlagMap{
		"--format": flag.New(
			"output format",
			scalar.String(
				scalar.Choices("table", "value-only"),
				scalar.Default("table"),
			),
			flag.Required(),
		),
	}
}

func widthFlag() flag.FlagMap {

	// TODO: figure out a good way to cache this for all width flags
	width := 0
	if term.IsTerminal(0) {
		termWidth, _, err := term.GetSize(0)
		if err == nil { // if there's not an error
			width = termWidth
		}
	}

	return flag.FlagMap{
		"--width": flag.New(
			"Width of the table. 0 means no limit",
			scalar.Int(
				scalar.Default(width),
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
	now := time.Now()
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
				scalar.Default(now),
			),
			flag.Required(),
		),
		"--update-time": flag.New(
			"Update time",
			scalar.New(
				datetime(),
				scalar.Default(now),
			),
			flag.Required(),
		),
	}
	return commonCreateFlags
}

func commonUpdateFlags() flag.FlagMap {

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
			scalar.String(),
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

// initEnvService reads the following from the flags:
//
// - --db-path
func initEnvService(ctx context.Context, passedFlags command.PassedFlags) (domain.EnvService, error) {
	sqliteDSN := passedFlags["--db-path"].(string)
	envService, err := sqlite.NewEnvService(ctx, sqliteDSN)
	if err != nil {
		//nolint:govet // we don't need the cancel if we err out
		return nil, fmt.Errorf("could not create env service: %w", err)
	}
	return envService, nil
}

type commonCreateArgs struct {
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
}

func mustGetCommonCreateArgs(pf command.PassedFlags) commonCreateArgs {
	return commonCreateArgs{
		Comment:    pf["--comment"].(string),
		CreateTime: pf["--create-time"].(time.Time),
		UpdateTime: pf["--update-time"].(time.Time),
	}
}

type commonUpdateArgs struct {
	Comment    *string
	CreateTime *time.Time
	NewName    *string
	UpdateTime *time.Time
}

func getCommonUpdateArgs(pf command.PassedFlags) commonUpdateArgs {
	return commonUpdateArgs{
		Comment:    ptrFromMap[string](pf, "--comment"),
		CreateTime: ptrFromMap[time.Time](pf, "--create-time"),
		NewName:    ptrFromMap[string](pf, "--new-name"),
		UpdateTime: ptrFromMap[time.Time](pf, "--update-time"),
	}
}

func mustGetConfirmArg(pf command.PassedFlags) bool {
	return pf["--confirm"].(bool)
}

func mustGetEnvNameArg(pf command.PassedFlags) string {
	return pf["--env-name"].(string)
}

func mustGetMaskArg(pf command.PassedFlags) bool {
	return pf["--mask"].(bool)
}

func mustGetNameArg(pf command.PassedFlags) string {
	return pf["--name"].(string)
}

func mustGetTimeoutArg(pf command.PassedFlags) time.Duration {
	return pf["--timeout"].(time.Duration)
}

func mustGetTimezoneArg(pf command.PassedFlags) string {
	return pf["--timezone"].(string)
}

func mustGetWidthArg(pf command.PassedFlags) int {
	return pf["--width"].(int)
}
