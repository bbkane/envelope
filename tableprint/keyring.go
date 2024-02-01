package tableprint

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/namedenv/domain"
)

func KeyringList(w io.Writer, keyringEntries []domain.KeyringEntry, errs []error, timezone Timezone) {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(w)

	if len(keyringEntries) > 0 {
		t.AppendHeader(table.Row{"Name", "Value", "Comment", "CreateTime", "UpdateTime"})
		for _, ke := range keyringEntries {
			t.AppendRow(table.Row{
				ke.Name,
				ke.Value,
				ke.Comment,
				formatTime(ke.CreateTime, timezone),
				formatTime(ke.UpdateTime, timezone),
			})
		}
		t.Render()
		t.ResetHeaders()
		t.ResetRows()
	} else {
		fmt.Fprintln(w, "no keyring entries found")
	}
	if len(errs) > 0 {
		t.AppendHeader(table.Row{"Error"})
		for _, e := range errs {
			t.AppendRow(table.Row{e})
		}
		t.Render()
		t.ResetHeaders()
		t.ResetRows()
	}
}
