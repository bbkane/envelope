package tableprint

import (
	"io"

	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(w io.Writer, envVar domain.EnvVar, timezone Timezone) {
	printKVTable(w, []kv{
		{"EnvName", envVar.EnvName},
		{"Name", envVar.Name},
		{"Value", envVar.Value},
		{"Comment", envVar.Comment},
		{"CreateTime", formatTime(envVar.CreateTime, timezone)},
		{"UpdateTime", formatTime(envVar.UpdateTime, timezone)},
	})
}
