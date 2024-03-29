package tableprint

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/envelope/domain"
)

func KeyringList(w io.Writer, keyringEntries []domain.KeyringEntry, errs []error, timezone Timezone) {
	if len(keyringEntries) > 0 {
		t := tableInit(w)
		for _, e := range keyringEntries {
			tableAddSection(t, []kv{
				{"Name", e.Name},
				{"Value", e.Value},
				{"Comment", e.Comment},
				{"CreateTime", formatTime(e.CreateTime, timezone)},
				{"UpdateTime", formatTime(e.UpdateTime, timezone)},
			})
		}
		t.Render()

	} else {
		fmt.Fprintln(w, "no keyring entries found")
	}

	if len(errs) > 0 {
		fmt.Fprintln(w, "Errors")
		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)
		t.SetOutputMirror(w)

		t.AppendHeader(table.Row{"Error"})
		for _, e := range errs {
			t.AppendRow(table.Row{e})
		}

	}
}
