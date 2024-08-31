package tableprint

import (
	"fmt"

	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(c CommonTablePrintArgs, envVar domain.EnvVar, envRefs []domain.EnvRef) {

	switch c.Format {
	case Format_Table:
		t := tableInit(c.W)
		tableAddSection(t, []row{
			newRow("EnvName", envVar.EnvName),
			newRow("Name", envVar.Name),
			newRow("Value", mask(c.Mask, envVar.Value)),
			newRow("Comment", envVar.Comment, skipRowIf(envVar.Comment == "")),
			newRow("CreateTime", formatTime(envVar.CreateTime, c.Tz)),
			newRow("UpdateTime", formatTime(envVar.UpdateTime, c.Tz)),
		})
		t.Render()

		if len(envRefs) > 0 {
			fmt.Fprintln(c.W, "EnvRefs")

			t := tableInit(c.W)

			for _, e := range envRefs {
				tableAddSection(t, []row{
					newRow("EnvName", e.EnvName),
					newRow("Name", e.Name),
					newRow("Comment", e.Comment, skipRowIf(e.Comment == "")),
				})
			}
			t.Render()
		}
	case Format_ValueOnly:
		fmt.Print(envVar.Value)
	default:
		panic("unexpected format: " + string(c.Format))
	}

}
