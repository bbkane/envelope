package tableprint

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(w io.Writer, envVar domain.EnvLocalVar, timezone Timezone) {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(w)

	//nolint:exhaustruct
	columnConfigs := []table.ColumnConfig{
		{Name: "Name"},
		{Name: "Value"},
	}

	t.SetColumnConfigs(columnConfigs)

	t.AppendHeader(table.Row{"Name", "Value"})
	t.AppendRows([]table.Row{
		{"EnvName", envVar.EnvName},
		{"Name", envVar.Name},
		{"Comment", envVar.Comment},
		{"CreateTime", formatTime(envVar.CreateTime, timezone)},
		{"UpdateTime", formatTime(envVar.UpdateTime, timezone)},
		{"Value", envVar.Value},
	})

	t.Render()
}
