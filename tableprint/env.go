package tableprint

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/namedenv/domain"
)

func EnvTable(w io.Writer, env domain.Env, timezone Timezone) {
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
		{"Name", env.Name},
		{"Comment", valOrEmpty(env.Comment)},
		{"CreateTime", formatTime(env.CreateTime, timezone)},
		{"UpdateTime", formatTime(env.UpdateTime, timezone)},
	})

	t.Render()
}
