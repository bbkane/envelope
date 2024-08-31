package tableprint

import (
	"fmt"

	"go.bbkane.com/envelope/domain"
)

func EnvRefShowPrint(c CommonTablePrintArgs, envRef domain.EnvRef, envVar domain.EnvVar) {

	switch c.Format {
	case Format_Table:
		t := tableInit(c.W)
		tableAddSection(t, []row{
			newRow("EnvName", envRef.EnvName),
			newRow("Name", envRef.Name),
			newRow("RefEnvName", envRef.RefEnvName),
			newRow("RefVarName", envRef.RevVarName),
			newRow("RefVarValue", mask(c.Mask, envVar.Value)),
			newRow("Comment", envRef.Comment, skipRowIf(envRef.Comment == "")),
			newRow("CreateTime", formatTime(envRef.CreateTime, c.Tz)),
			newRow("UpdateTime", formatTime(envRef.UpdateTime, c.Tz)),
		})
		t.Render()
	case Format_ValueOnly:
		fmt.Print(envVar.Value)
	default:
		panic("unexpected format: " + string(c.Format))
	}

}
