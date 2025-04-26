package tableprint

import (
	"fmt"

	"go.bbkane.com/envelope/models"
)

func VarRefShowPrint(c CommonTablePrintArgs, envRef models.VarRef, envVar models.Var) {

	switch c.Format {
	case Format_Table:
		t := newKeyValueTable(c.W, c.DesiredMaxWidth, len("RefVarValue"))
		createTime := formatTime(envRef.CreateTime, c.Tz)
		updateTime := formatTime(envRef.UpdateTime, c.Tz)
		t.Section(
			newRow("EnvName", envRef.EnvName),
			newRow("Name", envRef.Name),
			newRow("RefEnvName", envRef.RefEnvName),
			newRow("RefVarName", envRef.RevVarName),
			newRow("RefVarValue", mask(c.Mask, envVar.Value)),
			newRow("Comment", envRef.Comment, skipRowIf(envRef.Comment == "")),
			newRow("CreateTime", createTime),
			newRow("UpdateTime", updateTime, skipRowIf(envRef.CreateTime.Equal(envRef.UpdateTime))),
		)
		t.Render()
	case Format_ValueOnly:
		fmt.Print(envVar.Value)
	default:
		panic("unexpected format: " + string(c.Format))
	}

}
