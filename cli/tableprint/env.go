package tableprint

import (
	"fmt"
	"io"

	"go.bbkane.com/envelope/domain"
)

func EnvList(w io.Writer, envs []domain.Env, timezone Timezone) {

	if len(envs) > 0 {
		for _, e := range envs {
			printKVTable(w, []kv{
				{"Name", e.Name},
				{"Comment", e.Comment},
				{"CreateTime", formatTime(e.CreateTime, timezone)},
				{"UpdateTime", formatTime(e.UpdateTime, timezone)},
			})
		}
	} else {
		fmt.Fprintln(w, "no envs found")
	}
}

func EnvShowRun(
	w io.Writer,
	env domain.Env,
	localvars []domain.EnvVar,
	refs []domain.EnvRef,
	referencedVars []domain.EnvVar,
	timezone Timezone,
) {

	fmt.Fprintln(w, "Env")

	printKVTable(w, []kv{
		{"Name", env.Name},
		{"Comment", env.Comment},
		{"CreateTime", formatTime(env.CreateTime, timezone)},
		{"UpdateTime", formatTime(env.UpdateTime, timezone)},
	})

	if len(localvars) > 0 {
		fmt.Fprintln(w, "Vars")

		for _, e := range localvars {
			printKVTable(w, []kv{
				{"Name", e.Name},
				{"Value", e.Value},
				{"Comment", e.Comment},
			})
		}
	}

	if len(refs) > 0 {
		fmt.Fprintln(w, "Refs")

		for i := range len(refs) {
			printKVTable(w, []kv{
				{"Name", refs[i].Name},
				{"RefEnvName", referencedVars[i].EnvName},
				{"RefVarName", referencedVars[i].Name},
				{"RefVarValue", referencedVars[i].Value},
				{"Comment", refs[i].Comment},
			})
		}
	}

}
