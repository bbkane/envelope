package tableprint

import (
	"fmt"
	"io"

	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(w io.Writer, envVar domain.EnvVar, envRefs []domain.EnvRef, timezone Timezone) {
	printKVTable(w, []kv{
		{"EnvName", envVar.EnvName},
		{"Name", envVar.Name},
		{"Value", envVar.Value},
		{"Comment", envVar.Comment},
		{"CreateTime", formatTime(envVar.CreateTime, timezone)},
		{"UpdateTime", formatTime(envVar.UpdateTime, timezone)},
	})

	if len(envRefs) > 0 {
		fmt.Fprintln(w, "EnvRefs")

		for _, e := range envRefs {
			printKVTable(w, []kv{
				{"EnvName", e.EnvName},
				{"Name", e.Name},
				{"Comment", e.Comment},
			})
		}
	}
}
