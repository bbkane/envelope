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
			name:            "01_envCreate01",
			args:            envCreateTestCmd(dbName, envName01),
			expectActionErr: false,
		},
		{
			name:            "02_varCreate",
			args:            varCreateTestCmd(dbName, envName01, envVarName01, "val01"),
			expectActionErr: false,
		},
		{
			name:            "03_envCreate02",
			args:            envCreateTestCmd(dbName, envName02),
			expectActionErr: false,
		},
		{
			name: "04_varRefCreate",
			args: new(testCmdBuilder).Strs("var", "ref", "create").
				EnvName(envName02).Name(envRefName01).ZeroTimes().
				Strs("--ref-env-name", envName01).
				Strs("--ref-var-name", envVarName01).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "05_envShow02",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName02).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "06_varRefShow",
			args: new(testCmdBuilder).Strs("var", "ref", "show").
				EnvName(envName02).Name(envRefName01).Tz().Mask(false).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "07_envVarShow01",
			args: new(testCmdBuilder).Strs("var", "show").
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
			name:            "02_envLocalVarCreate",
			args:            varCreateTestCmd(dbName, envName01, envVarName01, "value"),
			expectActionErr: false,
		},
		{
			name: "03_envRefCreate",
			args: new(testCmdBuilder).Strs("var", "ref", "create").
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
			args: new(testCmdBuilder).Strs("var", "delete").
				EnvName(envName01).Name(envVarName01).Confirm(false).Finish(dbName),
			expectActionErr: true,
		},
		{
			name: "06_envRefDelete",
			args: new(testCmdBuilder).Strs("var", "ref", "delete").
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
