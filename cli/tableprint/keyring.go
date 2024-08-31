package tableprint

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/envelope/domain"
)

func KeyringList(w io.Writer, keyringEntries []domain.KeyringEntry, errs []error, timezone Timezone) {
	if len(keyringEntries) > 0 {
		t := NewKeyValueTable(w, 0, 0)
		for _, e := range keyringEntries {
			createTime := formatTime(e.CreateTime, timezone)
			updateTime := formatTime(e.UpdateTime, timezone)
			t.Section(
				newRow("Name", e.Name),
				newRow("Value", e.Value),
				newRow("Comment", e.Comment, skipRowIf(e.Comment == "")),
				newRow("CreateTime", createTime),
				newRow("UpdateTime", updateTime, skipRowIf(e.CreateTime == e.UpdateTime)),
			)
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
