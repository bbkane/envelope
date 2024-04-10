package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvLocalVarCreate(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbFile, err := os.CreateTemp(os.TempDir(), "envelope-test-")
	require.NoError(t, err)
	err = dbFile.Close()
	require.NoError(t, err)

	t.Log("dbFile:", dbFile.Name())

	tests := []testcase{
		{
			name:            "01_envCreate",
			args:            createEnv(dbFile.Name(), envName01),
			expectActionErr: false,
		},
		{
			name: "02_envLocalVarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).Strs("--value", "value").
				ZeroTimes().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "03_envLocalVarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").
				EnvName(envName01).Name(envVarName01).Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
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

	tests := []testcase{

		{
			name:            "01_envCreate",
			args:            createEnv(dbFile.Name(), envName01),
			expectActionErr: false,
		},
		{
			name: "02_envLocalVarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).
				Strs("--value", "value").ZeroTimes().Finish(dbFile.Name()),

			expectActionErr: false,
		},
		{
			name: "03_envLocalVarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").EnvName(envName01).
				Name(envVarName01).Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "04_envLocalVarDelete",
			args: new(testCmdBuilder).Strs("env", "var", "delete").Strs("--confirm", "false").
				EnvName(envName01).Name(envVarName01).Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "05_envLocalVarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").EnvName(envName01).
				Name(envVarName01).Tz().Finish(dbFile.Name()),
			expectActionErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}

func TestEnvNonUniqueNames(t *testing.T) {
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
			name:            "01_envCreate",
			args:            createEnv(dbFile.Name(), envName01),
			expectActionErr: false,
		},
		{
			name: "02_envLocalVarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).Strs("--value", "value").
				ZeroTimes().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "03_envRefCreateSameName",
			args: new(testCmdBuilder).Strs("env", "ref", "create").
				EnvName(envName01).Name("key").Strs("--ref-env-name", envName01).
				Strs("--ref-var-name", envVarName01).ZeroTimes().Finish(dbFile.Name()),
			expectActionErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}
