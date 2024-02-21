package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
)

func TestEnvLocalVarCreate(t *testing.T) {
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
			name: "02_envLocalVarCreate",
			args: []string{
				"envelope", "env", "var", "create",
				"--db-path", dbFile.Name(),
				"--env-name", "env_name",
				"--name", "key",
				"--value", "value",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "03_envLocalVarShow",
			args: []string{
				"envelope", "env", "var", "show",
				"--db-path", dbFile.Name(),
				"--env-name", "env_name",
				"--name", "key",
				"--timezone", "utc",
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

func TestEnvLocalVarDelete(t *testing.T) {
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
			name: "02_envLocalVarCreate",
			args: []string{
				"envelope", "env", "var", "create",
				"--db-path", dbFile.Name(),
				"--env-name", "env_name",
				"--name", "key",
				"--value", "value",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "03_envLocalVarShow",
			args: []string{
				"envelope", "env", "var", "show",
				"--db-path", dbFile.Name(),
				"--env-name", "env_name",
				"--name", "key",
				"--timezone", "utc",
			},
			expectActionErr: false,
		},
		{
			name: "04_envLocalVarDelete",
			args: []string{
				"envelope", "env", "var", "delete",
				"--confirm", "false",
				"--db-path", dbFile.Name(),
				"--env-name", "env_name",
				"--name", "key",
			},
			expectActionErr: false,
		},
		{
			name: "05_envLocalVarShow",
			args: []string{
				"envelope", "env", "var", "show",
				"--db-path", dbFile.Name(),
				"--env-name", "env_name",
				"--name", "key",
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
