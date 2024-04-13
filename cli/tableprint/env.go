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
	c CommonTablePrintArgs,
	env domain.Env,
	localvars []domain.EnvVar,
	refs []domain.EnvRef,
	referencedVars []domain.EnvVar,
) {

	fmt.Fprintln(c.W, "Env")

	t := tableInit(c.W)
	tableAddSection(t, []kv{
		{"Name", env.Name},
		{"Comment", env.Comment},
		{"CreateTime", formatTime(env.CreateTime, c.Tz)},
		{"UpdateTime", formatTime(env.UpdateTime, c.Tz)},
	})
	t.Render()

	if len(localvars) > 0 {
		fmt.Fprintln(c.W, "Vars")

		t := tableInit(c.W)
		for _, e := range localvars {
			tableAddSection(t, []kv{
				{"Name", e.Name},
				{"Value", mask(c.Mask, e.Value)},
				{"Comment", e.Comment},
			})
		}
		t.Render()
	}

	if len(refs) > 0 {
		fmt.Fprintln(c.W, "Refs")
		t := tableInit(c.W)

		for i := range len(refs) {
			tableAddSection(t, []kv{
				{"Name", refs[i].Name},
				{"RefEnvName", referencedVars[i].EnvName},
				{"RefVarName", referencedVars[i].Name},
				{"RefVarValue", mask(c.Mask, referencedVars[i].Value)},
				{"Comment", refs[i].Comment},
			})
		}
		t.Render()

	}

}
