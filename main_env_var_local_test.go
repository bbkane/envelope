package main

import (
	"os"
	"testing"
)

func TestEnvLocalVarCreate(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbName := createTempDB(t)

	tests := []testcase{
		{
			name:            "01_envCreate",
			args:            createEnv(dbName, envName01),
			expectActionErr: false,
		},
		{
			name: "02_envLocalVarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).Strs("--value", "value").
				ZeroTimes().Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "03_envLocalVarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").
				EnvName(envName01).Name(envVarName01).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Mask(false).Finish(dbName),
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

	dbName := createTempDB(t)

	tests := []testcase{

		{
			name:            "01_envCreate",
			args:            createEnv(dbName, envName01),
			expectActionErr: false,
		},
		{
			name: "02_envLocalVarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).
				Strs("--value", "value").ZeroTimes().Finish(dbName),

			expectActionErr: false,
		},
		{
			name: "03_envLocalVarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").EnvName(envName01).
				Name(envVarName01).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "04_envLocalVarDelete",
			args: new(testCmdBuilder).Strs("env", "var", "delete").Confirm(false).
				EnvName(envName01).Name(envVarName01).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "05_envLocalVarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").EnvName(envName01).
				Name(envVarName01).Tz().Mask(false).Finish(dbName),
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

	dbName := createTempDB(t)

	tests := []testcase{
		{
			name:            "01_envCreate",
			args:            createEnv(dbName, envName01),
			expectActionErr: false,
		},
		{
			name: "02_envLocalVarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).Strs("--value", "value").
				ZeroTimes().Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "03_envRefCreateSameName",
			args: new(testCmdBuilder).Strs("env", "ref", "create").
				EnvName(envName01).Name(envVarName01).Strs("--ref-env-name", envName01).
				Strs("--ref-var-name", envVarName01).ZeroTimes().Finish(dbName),
			expectActionErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}
