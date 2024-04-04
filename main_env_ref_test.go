package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
)

func TestEnvRefCreate(t *testing.T) {
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
			name: "01_env01Create",
			args: []string{
				"envelope", "env", "create",
				"--db-path", dbFile.Name(),
				"--name", "env01",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "02_env01VarCreate",
			args: []string{
				"envelope", "env", "var", "create",
				"--db-path", dbFile.Name(),
				"--env-name", "env01",
				"--name", "var01",
				"--value", "val01",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "03_env02Create",
			args: []string{
				"envelope", "env", "create",
				"--db-path", dbFile.Name(),
				"--name", "env02",
				"--create-time", zeroTime,
				"--update-time", zeroTime,
			},
			expectActionErr: false,
		},
		{
			name: "04_env02RefCreate",
			args: []string{
				"envelope", "env", "ref", "create",
				"--db-path", dbFile.Name(),
				"--env-name", "env02",
				"--name", "ref01",
				"--ref-env-name", "env01",
				"--ref-var-name", "var01",
			},
			expectActionErr: false,
		},
		{
			name: "05_env02Show",
			args: []string{
				"envelope", "env", "show",
				"--db-path", dbFile.Name(),
				"--name", "env02",
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
