package tableprint

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/namedenv/domain"
)

func EnvShowRun(w io.Writer, env domain.Env, localvars []domain.EnvLocalVar, timezone Timezone) {

	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(w)

	fmt.Fprintln(w, "Env")

	t.AppendHeader(table.Row{"Name", "Value"})
	t.AppendRows([]table.Row{
		{"Name", env.Name},
		{"Comment", env.Comment},
		{"CreateTime", formatTime(env.CreateTime, timezone)},
		{"UpdateTime", formatTime(env.UpdateTime, timezone)},
	})

	t.Render()

	t.ResetHeaders()
	t.ResetRows()

	if len(localvars) > 0 {
		fmt.Fprintln(w, "LocalVars")

		t.AppendHeader(table.Row{"Name", "Value", "Comment", "CreateTime", "UpdateTime"})
		for _, e := range localvars {
			t.AppendRow(table.Row{
				e.Name, e.Value, e.Comment, e.CreateTime, e.UpdateTime,
			})
		}

		t.Render()

		t.ResetHeaders()
		t.ResetRows()
	}

}
