package tableprint

import (
	"fmt"

	"go.bbkane.com/envelope/domain"
)

func EnvRefShowPrint(c CommonTablePrintArgs, envRef domain.EnvRef, envVar domain.EnvVar) {

	switch c.Format {
	case Format_Table:
		t := tableInit(c.W)
		tableAddSection(t, []kv{
			{"EnvName", envRef.EnvName},
			{"Name", envRef.Name},
			{"RefEnvName", envRef.RefEnvName},
			{"RefVarName", envRef.RevVarName},
			{"RefVarValue", mask(c.Mask, envVar.Value)},
			{"Comment", envRef.Comment},
			{"CreateTime", formatTime(envRef.CreateTime, c.Tz)},
			{"UpdateTime", formatTime(envRef.UpdateTime, c.Tz)},
		})
		t.Render()
	case Format_ValueOnly:
		fmt.Print(envVar.Value)
	default:
		panic("unexpected format: " + string(c.Format))
	}

}
