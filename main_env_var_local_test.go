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
			args:            envCreateTestCmd(dbName, envName01),
			expectActionErr: false,
		},
		{
			name:            "02_envLocalVarCreate",
			args:            envVarCreateTestCmd(dbName, envName01, envVarName01, "value"),
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
			args:            envCreateTestCmd(dbName, envName01),
			expectActionErr: false,
		},
		{
			name:            "02_envLocalVarCreate",
			args:            envVarCreateTestCmd(dbName, envName01, envVarName01, "value"),
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
			args:            envCreateTestCmd(dbName, envName01),
			expectActionErr: false,
		},
		{
			name:            "02_envLocalVarCreate",
			args:            envVarCreateTestCmd(dbName, envName01, envVarName01, "value"),
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

func TestEnvVarUpdate(t *testing.T) {
	t.Parallel()
	updateGolden := os.Getenv("ENVELOPE_TEST_UPDATE_GOLDEN") != ""

	dbName := createTempDB(t)

	const e1 = "e1"
	const e2 = "e2"
	const v1 = "v1"
	const v2 = "v2"
	const v3 = "v3"

	tests := []testcase{
		// setup
		{
			name:            "01_e1Create",
			args:            envCreateTestCmd(dbName, e1),
			expectActionErr: false,
		},
		{
			name:            "02_e2Create",
			args:            envCreateTestCmd(dbName, e2),
			expectActionErr: false,
		},
		{
			name:            "03_e1v1Create",
			args:            envVarCreateTestCmd(dbName, e1, v1, "val"),
			expectActionErr: false,
		},
		{
			name:            "04_e1v2Create",
			args:            envVarCreateTestCmd(dbName, e1, v2, "val"),
			expectActionErr: false,
		},
		{
			name:            "05_e2v2Create",
			args:            envVarCreateTestCmd(dbName, e2, v2, "val"),
			expectActionErr: false,
		},
		{
			name:            "06_e2v3Create",
			args:            envVarCreateTestCmd(dbName, e2, v3, "val"),
			expectActionErr: false,
		},
		// test updates
		{
			name:            "07_nameUpdateConflict",
			args:            new(testCmdBuilder).Strs("env", "var", "update").Confirm(false).UpdateTime(oneTime).EnvName(e1).Name(v1).Strs("--new-name", v2).Finish(dbName),
			expectActionErr: true,
		},
		{
			name:            "08_envUpdateConflict",
			args:            new(testCmdBuilder).Strs("env", "var", "update").Confirm(false).UpdateTime(oneTime).EnvName(e1).Name(v2).Strs("--new-env-name", e2).Finish(dbName),
			expectActionErr: true,
		},
		{
			name:            "09_envUpdateNameUpdateConflict",
			args:            new(testCmdBuilder).Strs("env", "var", "update").Confirm(false).UpdateTime(oneTime).EnvName(e1).Name(v1).Strs("--new-env-name", e2).Strs("--new-name", v3).Finish(dbName),
			expectActionErr: true,
		},
		{
			name:            "10_emptyUpdate",
			args:            new(testCmdBuilder).Strs("env", "var", "update").Confirm(false).UpdateTime(oneTime).EnvName(e1).Name(v1).Strs("--new-env-name", e1).Strs("--new-name", v1).Finish(dbName),
			expectActionErr: false,
		},
		{
			name: "11_updateOtherFields",
			args: new(testCmdBuilder).Strs("env", "var", "update").
				Confirm(false).UpdateTime(oneTime).CreateTime(oneTime).EnvName(e1).Name(v1).
				Strs("--new-env-name", e2, "--new-name", "NEW", "--comment", "NEW", "--value", "NEW").Finish(dbName),
			expectActionErr: false,
		},
		{
			name:            "12_e1Show",
			args:            envShowTestCmd(dbName, e1),
			expectActionErr: false,
		},
		{
			name:            "13_e2Show",
			args:            envShowTestCmd(dbName, e2),
			expectActionErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, tt, updateGolden)
		})
	}
}
