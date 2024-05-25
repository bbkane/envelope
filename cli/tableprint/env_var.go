package tableprint

import (
	"fmt"

	"go.bbkane.com/envelope/domain"
)

func EnvLocalVarShowPrint(c CommonTablePrintArgs, envVar domain.EnvVar, envRefs []domain.EnvRef) {

	switch c.Format {
	case Format_Table:
		t := tableInit(c.W)
		tableAddSection(t, []kv{
			{"EnvName", envVar.EnvName},
			{"Name", envVar.Name},
			{"Value", mask(c.Mask, envVar.Value)},
			{"Comment", envVar.Comment},
			{"CreateTime", formatTime(envVar.CreateTime, c.Tz)},
			{"UpdateTime", formatTime(envVar.UpdateTime, c.Tz)},
		})
		t.Render()

		if len(envRefs) > 0 {
			fmt.Fprintln(c.W, "EnvRefs")

			t := tableInit(c.W)

			for _, e := range envRefs {
				tableAddSection(t, []kv{
					{"EnvName", e.EnvName},
					{"Name", e.Name},
					{"Comment", e.Comment},
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
