package main

import (
	"testing"

	"go.bbkane.com/warg"
)

// abstract some common test functionality to make writing
// tests shorter

const zeroTime = "0001-01-01T00:00:00Z"

const envName01 = "env_name"
const envVarName01 = "key"

// const envRefName01 = "envRefName01"

func createEnv(dbPath string, envName string) []string {
	return new(testCmdBuilder).
		Strs("env", "create").
		Name(envName).
		ZeroTimes().
		Finish(dbPath)
}

type testCmdBuilder struct {
	cmd []string
}

func (tcb *testCmdBuilder) Strs(args ...string) *testCmdBuilder {
	tcb.cmd = append(tcb.cmd, args...)
	return tcb
}

func (tcb *testCmdBuilder) Finish(dbPath string) []string {
	cmd := []string{"envelope", "--db-path", dbPath}
	cmd = append(cmd, tcb.cmd...)
	return cmd
}

func (tcb *testCmdBuilder) Name(name string) *testCmdBuilder {
	return tcb.Strs("--name", name)
}

func (tcb *testCmdBuilder) ZeroTimes() *testCmdBuilder {
	return tcb.Strs("--create-time", zeroTime, "--update-time", zeroTime)
}

func (tcb *testCmdBuilder) EnvName(envName string) *testCmdBuilder {
	return tcb.Strs("--env-name", envName)
}

// Tz adds "--timezone", "utc"
func (tcb *testCmdBuilder) Tz() *testCmdBuilder {
	return tcb.Strs("--timezone", "utc")
}

type testcase struct {
	name            string
	args            []string
	expectActionErr bool
}

func goldenTest(t *testing.T, tt testcase, updateGolden bool) {
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
}
