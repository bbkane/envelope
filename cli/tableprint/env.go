package tableprint

import (
	"fmt"
	"io"

	"go.bbkane.com/envelope/domain"
)

func EnvList(w io.Writer, envs []domain.Env, timezone Timezone) {

	if len(envs) > 0 {
		t := tableInit(w)
		for _, e := range envs {
			tableAddSection(t, []kv{
				{"Name", e.Name},
				{"Comment", e.Comment},
				{"CreateTime", formatTime(e.CreateTime, timezone)},
				{"UpdateTime", formatTime(e.UpdateTime, timezone)},
			})
		}
		t.Render()
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

	t := tableInit(w)
	tableAddSection(t, []kv{
		{"Name", env.Name},
		{"Comment", env.Comment},
		{"CreateTime", formatTime(env.CreateTime, timezone)},
		{"UpdateTime", formatTime(env.UpdateTime, timezone)},
	})
	t.Render()

	if len(localvars) > 0 {
		fmt.Fprintln(w, "Vars")

		t := tableInit(w)
		for _, e := range localvars {
			tableAddSection(t, []kv{
				{"Name", e.Name},
				{"Value", e.Value},
				{"Comment", e.Comment},
			})
		}
		t.Render()
	}

	if len(refs) > 0 {
		fmt.Fprintln(w, "Refs")
		t := tableInit(w)

		for i := range len(refs) {
			tableAddSection(t, []kv{
				{"Name", refs[i].Name},
				{"RefEnvName", referencedVars[i].EnvName},
				{"RefVarName", referencedVars[i].Name},
				{"RefVarValue", referencedVars[i].Value},
				{"Comment", refs[i].Comment},
			})
		}
		t.Render()

	}

}
