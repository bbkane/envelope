package tableprint

import (
	"fmt"
	"io"

	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(w io.Writer, envVar domain.EnvVar, envRefs []domain.EnvRef, timezone Timezone) {
	t := tableInit(w)
	tableAddSection(t, []kv{
		{"EnvName", envVar.EnvName},
		{"Name", envVar.Name},
		{"Value", envVar.Value},
		{"Comment", envVar.Comment},
		{"CreateTime", formatTime(envVar.CreateTime, timezone)},
		{"UpdateTime", formatTime(envVar.UpdateTime, timezone)},
	})
	t.Render()

	if len(envRefs) > 0 {
		fmt.Fprintln(w, "EnvRefs")

		t := tableInit(w)

		for _, e := range envRefs {
			tableAddSection(t, []kv{
				{"EnvName", e.EnvName},
				{"Name", e.Name},
				{"Comment", e.Comment},
			})
		}
		t.Render()
	}
}
