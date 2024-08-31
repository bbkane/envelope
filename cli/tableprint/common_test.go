package tableprint

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableTruncation(t *testing.T) {
	t.Parallel()
	// ╭─────┬──────╮
	// │ key │ v... │
	// ╰─────┴──────╯
	// 12345678901234

	expected := `╭─────┬──────╮
│ key │ v... │
╰─────┴──────╯
`

	buf := new(bytes.Buffer)
	tw := newKeyValueTable(buf, 14, len("key"))
	tw.Section(
		newRow("key", "value"),
	)
	tw.Render()
	actual := buf.String()

	require.Equal(t, expected, actual)
}

func TestTableTruncationNoTruncate(t *testing.T) {
	t.Parallel()
	// ╭─────┬───────╮
	// │ key │ value │
	// ╰─────┴───────╯
	// 123456789012345

	expected := `╭─────┬───────╮
│ key │ value │
╰─────┴───────╯
`

	buf := new(bytes.Buffer)
	tw := newKeyValueTable(buf, 15, len("key"))
	tw.Section(
		newRow("key", "value"),
	)
	tw.Render()
	actual := buf.String()

	require.Equal(t, expected, actual)
}

func TestTableTruncationToNarrow(t *testing.T) {
	t.Parallel()
	// ╭─────┬───────╮
	// │ key │ value │
	// ╰─────┴───────╯
	// 123456789012345

	expected := `╭─────┬───────╮
│ key │ value │
╰─────┴───────╯
`

	buf := new(bytes.Buffer)
	tw := newKeyValueTable(buf, 5, len("key"))
	tw.Section(
		newRow("key", "value"),
	)
	tw.Render()
	actual := buf.String()

	require.Equal(t, expected, actual)
}
