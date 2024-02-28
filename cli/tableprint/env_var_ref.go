package tableprint

import (
	"io"

	"go.bbkane.com/envelope/domain"
)

func EnvRefShowPrint(w io.Writer, envRef domain.EnvRef, envVar domain.EnvVar, timezone Timezone) {
	t := tableInit(w)
	tableAddSection(t, []kv{
		{"EnvName", envRef.EnvName},
		{"Name", envRef.Name},
		{"RefEnvName", envRef.RefEnvName},
		{"RefVarName", envRef.RevVarName},
		{"RefVarValue", envVar.Value},
		{"Comment", envRef.Comment},
		{"CreateTime", formatTime(envRef.CreateTime, timezone)},
		{"UpdateTime", formatTime(envRef.UpdateTime, timezone)},
	})
	t.Render()
}
