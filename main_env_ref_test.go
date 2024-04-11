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
			args:            createEnv(dbName, envName01),
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
			args:            createEnv(dbName, envName02),
			expectActionErr: false,
		},
		{
			name: "04_env02RefCreate",
			args: new(testCmdBuilder).Strs("env", "ref", "create").
				EnvName(envName02).Name(envRefName01).Strs("--ref-env-name", envName01).
				Strs("--ref-var-name", envVarName01).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "05_env02Show",
			args: new(testCmdBuilder).Strs("env", "show").
				Name(envName02).Tz().Finish(dbName),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}
