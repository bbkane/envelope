package tableprint

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.bbkane.com/envelope/domain"
)

func KeyringList(c CommonTablePrintArgs, keyringEntries []domain.KeyringEntry, errs []error) {
	if len(keyringEntries) > 0 {
		t := newKeyValueTable(c.W, c.DesiredMaxWidth, len("CreateTime"))
		for _, e := range keyringEntries {
			createTime := formatTime(e.CreateTime, c.Tz)
			updateTime := formatTime(e.UpdateTime, c.Tz)
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
		fmt.Fprintln(c.W, "no keyring entries found")
	}

	if len(errs) > 0 {
		fmt.Fprintln(c.W, "Errors")
		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)
		t.SetOutputMirror(c.W)

		t.AppendHeader(table.Row{"Error"})
		for _, e := range errs {
			t.AppendRow(table.Row{e})
		}

	}
}
