package main

import (
	"os"
	"testing"
)

func TestEnvRefCreate(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbName := createTempDB(t)

	t.Log("dbFile:", dbName)

	tests := []testcase{
		{
			name:            "01_env01Create",
			args:            envCreateTestCmd(dbName, envName01),
			expectActionErr: false,
		},
		{
			name: "02_env01VarCreate",
			args: new(testCmdBuilder).Strs("env", "var", "create").
				EnvName(envName01).Name(envVarName01).Strs("--value", "val01").
				ZeroTimes().Finish(dbName),
			expectActionErr: false,
		},
		{
			name:            "03_env02Create",
			args:            envCreateTestCmd(dbName, envName02),
			expectActionErr: false,
		},
		{
			name: "04_env02RefCreate",
			args: new(testCmdBuilder).Strs("env", "ref", "create").
				EnvName(envName02).Name(envRefName01).ZeroTimes().
				Strs("--ref-env-name", envName01).
				Strs("--ref-var-name", envVarName01).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "05_env02Show",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName02).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "06_env02RefShow",
			args: new(testCmdBuilder).Strs("env", "ref", "show").
				EnvName(envName02).Name(envRefName01).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "07_env01VarShow",
			args: new(testCmdBuilder).Strs("env", "var", "show").
				EnvName(envName01).Name(envVarName01).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}

func TestEnvRefDelete(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbName := createTempDB(t)

	tests := []testcase{
		{
			name:            "01_envCreate",
			args:            envCreateTestCmd(dbName, envName01),
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
			name: "03_envRefCreate",
			args: new(testCmdBuilder).Strs("env", "ref", "create").
				EnvName(envName01).Name(envRefName01).Strs("--ref-env-name", envName01).
				Strs("--ref-var-name", envVarName01).ZeroTimes().Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "04_envShow",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName01).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "05_envVarDeleteRestricted",
			args: new(testCmdBuilder).Strs("env", "var", "delete").
				EnvName(envName01).Name(envVarName01).Confirm(false).Finish(dbName),
			expectActionErr: true,
		},
		{
			name: "06_envRefDelete",
			args: new(testCmdBuilder).Strs("env", "ref", "delete").
				EnvName(envName01).Name(envRefName01).Confirm(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "07_envShow",
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
