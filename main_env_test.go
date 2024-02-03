package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
)

func TestBuildApp(t *testing.T) {
	app := buildApp()

	if err := app.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvCreate(t *testing.T) {
	updateGolden := os.Getenv("envelope_TEST_UPDATE_GOLDEN") != ""

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
				"--sqlite-dsn", dbFile.Name(),
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
				"--sqlite-dsn", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "03_envList",
			args: []string{
				"envelope", "env", "list",
				"--sqlite-dsn", dbFile.Name(),
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
					ExpectActionErr: false,
				},
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(warg.LookupMap(nil)),
			)
		})
	}
}

func TestEnvDelete(t *testing.T) {
	updateGolden := os.Getenv("envelope_TEST_UPDATE_GOLDEN") != ""

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
				"--sqlite-dsn", dbFile.Name(),
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
				"--sqlite-dsn", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "03_envDelete",
			args: []string{
				"envelope", "env", "delete",
				"--sqlite-dsn", dbFile.Name(),
				"--name", "env_name",
			},
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: []string{
				"envelope", "env", "show",
				"--sqlite-dsn", dbFile.Name(),
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
	updateGolden := os.Getenv("envelope_TEST_UPDATE_GOLDEN") != ""

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
				"--sqlite-dsn", dbFile.Name(),
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
				"--sqlite-dsn", dbFile.Name(),
				"--name", "env_name",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "03_envUpdate",
			args: []string{
				"envelope", "env", "update",
				"--sqlite-dsn", dbFile.Name(),
				"--name", "env_name",
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
				"--sqlite-dsn", dbFile.Name(),
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
					ExpectActionErr: false,
				},
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(warg.LookupMap(nil)),
			)
		})
	}
}