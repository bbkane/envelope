package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

	tests := []testcase{
		{
			name:            "01_envCreate",
			args:            createEnv(dbFile.Name(), envName01),
			expectActionErr: false,
		},
		{
			name: "02_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "03_envList",
			args: new(testCmdBuilder).Strs("env", "list").
				Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
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

	tests := []testcase{
		{
			name:            "01_envCreate",
			args:            createEnv(dbFile.Name(), envName01),
			expectActionErr: false,
		},
		{
			name: "02_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "03_envDelete",
			args: new(testCmdBuilder).Strs("env", "delete").
				Confirm(false).Name(envName01).Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Finish(dbFile.Name()),
			expectActionErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
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

	tests := []testcase{
		{
			name: "01_envCreate",
			args: new(testCmdBuilder).Strs("env", "create").
				Name(envName01).ZeroTimes().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "02_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "03_envUpdate",
			args: new(testCmdBuilder).Strs("env", "update").
				Name(envName01).Confirm(false).Strs("--comment", "a comment").
				Strs("--create-time", oneTime).Strs("--new-name", "new_name").
				Strs("--update-time", oneTime).Finish(dbFile.Name()),
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name("new_name").Tz().Finish(dbFile.Name()),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
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

	tests := []testcase{
		{
			name: "01_envPrintScript",
			args: new(testCmdBuilder).Strs("env", "print-script").
				Name("non-existent-env").Finish(dbFile.Name()),
			expectActionErr: true,
		},
		{
			name: "01_envPrintScriptNoEnvNoProblem",
			args: new(testCmdBuilder).Strs("env", "print-script").
				Name("non-existent-env").Strs("--no-env-no-problem", "true").
				Finish(dbFile.Name()),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}
