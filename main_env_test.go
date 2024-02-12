package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
)

func TestBuildApp(t *testing.T) {
	t.Parallel()
	app := buildApp()

	if err := app.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvCreate(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbFile, err := os.CreateTemp(os.TempDir(), "envelope-test-")
	require.NoError(t, err)
	err = dbFile.Close()
	require.NoError(t, err)

	t.Log("dbFile:", dbFile.Name())

	const zeroTime = "0001-01-01T00:00:00Z"

	tests := []struct {
		name            string
		args            []string
		expectActionErr bool
	}{
		{
			name: "01_envCreate",
			args: []string{
				"envelope", "env", "create",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "02_envShow",
			args: []string{
				"envelope", "env", "show",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "03_envList",
			args: []string{
				"envelope", "env", "list",
				"--db-path", dbFile.Name(),
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				warg.GoldenTestArgs{
					App:             buildApp(),
					UpdateGolden:    updateGolden,
					ExpectActionErr: tt.expectActionErr,
				},
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(warg.LookupMap(nil)),
			)
		})
	}
}

func TestEnvDelete(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbFile, err := os.CreateTemp(os.TempDir(), "envelope-test-")
	require.NoError(t, err)
	err = dbFile.Close()
	require.NoError(t, err)

	t.Log("dbFile:", dbFile.Name())

	const zeroTime = "0001-01-01T00:00:00Z"

	tests := []struct {
		name            string
		args            []string
		expectActionErr bool
	}{
		{
			name: "01_envCreate",
			args: []string{
				"envelope", "env", "create",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "02_envShow",
			args: []string{
				"envelope", "env", "show",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "03_envDelete",
			args: []string{
				"envelope", "env", "delete",
				"--db-path", dbFile.Name(),
				"--confirm", "false",
				"--name", "env_name",
			},
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: []string{
				"envelope", "env", "show",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				warg.GoldenTestArgs{
					App:             buildApp(),
					UpdateGolden:    updateGolden,
					ExpectActionErr: tt.expectActionErr,
				},
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(warg.LookupMap(nil)),
			)
		})
	}
}

func TestEnvUpdate(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbFile, err := os.CreateTemp(os.TempDir(), "envelope-test-")
	require.NoError(t, err)
	err = dbFile.Close()
	require.NoError(t, err)

	t.Log("dbFile:", dbFile.Name())

	const zeroTime = "0001-01-01T00:00:00Z"
	const oneTime = "0001-01-01T01:00:00Z"

	tests := []struct {
		name            string
		args            []string
		expectActionErr bool
	}{
		{
			name: "01_envCreate",
			args: []string{
				"envelope", "env", "create",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "02_envShow",
			args: []string{
				"envelope", "env", "show",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "03_envUpdate",
			args: []string{
				"envelope", "env", "update",
				"--db-path", dbFile.Name(),
				"--name", "env_name",
				"--confirm", "false",
				"--comment", "a comment",
				"--create-time", oneTime,
				"--new-name", "new_name",
				"--update-time", oneTime,
			},
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: []string{
				"envelope", "env", "show",
				"--db-path", dbFile.Name(),
				"--name", "new_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				warg.GoldenTestArgs{
					App:             buildApp(),
					UpdateGolden:    updateGolden,
					ExpectActionErr: tt.expectActionErr,
				},
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(warg.LookupMap(nil)),
			)
		})
	}
}

func TestEnvPrintScript(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbFile, err := os.CreateTemp(os.TempDir(), "envelope-test-")
	require.NoError(t, err)
	err = dbFile.Close()
	require.NoError(t, err)

	t.Log("dbFile:", dbFile.Name())

	tests := []struct {
		name            string
		args            []string
		expectActionErr bool
	}{
		{
			name: "01_envPrintScript",
			args: []string{
				"envelope", "env", "print-script",
				"--db-path", dbFile.Name(),
				"--name", "non-existent-env",
			},
			expectActionErr: true,
		},
		{
			name: "01_envPrintScriptNoEnvNoProblem",
			args: []string{
				"envelope", "env", "print-script",
				"--db-path", dbFile.Name(),
				"--name", "non-existent-env",
				"--no-env-no-problem", "true",
			},
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				warg.GoldenTestArgs{
					App:             buildApp(),
					UpdateGolden:    updateGolden,
					ExpectActionErr: tt.expectActionErr,
				},
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(warg.LookupMap(nil)),
			)
		})
	}
}
