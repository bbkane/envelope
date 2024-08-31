package tableprint

import (
	"fmt"

	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(c CommonTablePrintArgs, envVar domain.EnvVar, envRefs []domain.EnvRef) {

	switch c.Format {
	case Format_Table:
		t := newKeyValueTable(c.W, c.DesiredMaxWidth, len("CreateTime"))
		createTime := formatTime(envVar.CreateTime, c.Tz)
		updateTime := formatTime(envVar.UpdateTime, c.Tz)
		t.Section(
			newRow("EnvName", envVar.EnvName),
			newRow("Name", envVar.Name),
			newRow("Value", mask(c.Mask, envVar.Value)),
			newRow("Comment", envVar.Comment, skipRowIf(envVar.Comment == "")),
			newRow("CreateTime", createTime),
			newRow("UpdateTime", updateTime, skipRowIf(envVar.CreateTime == envVar.UpdateTime)),
		)
		t.Render()

		if len(envRefs) > 0 {
			fmt.Fprintln(c.W, "EnvRefs")

			t := newKeyValueTable(c.W, c.DesiredMaxWidth, 0)
			for _, e := range envRefs {
				t.Section(
					newRow("EnvName", e.EnvName),
					newRow("Name", e.Name),
					newRow("Comment", e.Comment, skipRowIf(e.Comment == "")),
				)
			}
			t.Render()
		}
	case Format_ValueOnly:
		fmt.Print(envVar.Value)
	default:
		panic("unexpected format: " + string(c.Format))
	}

}
