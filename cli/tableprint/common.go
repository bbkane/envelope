package tableprint

import (
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Timezone string

const (
	Timezone_Local = "local"
	Timezone_UTC   = "utc"
)

// Format for CLI output
type Format string

const (
	Format_Table     = "table"
	Format_ValueOnly = "value-only"
)

type CommonTablePrintArgs struct {
	Format Format
	Mask   bool
	Tz     Timezone
	W      io.Writer
}

func mask(mask bool, val string) string {
	if mask {
		if len(val) < 2 {
			return "**"
		} else {
			return val[:2] + "****"
		}
	}
	return val
}

// truncate truncates a string to max-3 characters and appends "..." if the string is longer than max. As a special case, if max is 0, the original string is returned.
func truncate(s string, max int) string {
	if max == 0 || len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatTime(t time.Time, timezone Timezone) string {
	timeFormat := "Mon 2006-01-02"
	switch timezone {
	case Timezone_Local:
		return t.Local().Format(timeFormat)
	case Timezone_UTC:
		return t.UTC().Format(timeFormat)
	default:
		panic("unknown timezone: " + timezone)
	}
}

type row struct {
	Key   string
	Value string
	Skip  bool
}

type rowOpt func(*row)

func skipRowIf(skip bool) rowOpt {
	return func(r *row) {
		r.Skip = skip
	}
}

func newRow(key string, value string, opts ...rowOpt) row {

	r := row{
		Key:   key,
		Value: value,
		Skip:  false,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

type keyValueTable struct {
	t               table.Writer
	truncationWidth int
}

// newKeyValueTable creates a new table and tries to fit it into desiredMaxWidth
// desiredMaxWidth is ignored if it == 0, or if it is less than the minimum width possible
//
//	width = len(key) + len(truncated_value) + len(padding)
//
// if desiredMaxWidth < len(key) + len(truncated_value) + len(padding) , it is ignored.
func newKeyValueTable(w io.Writer, desiredMaxWidth int, maxKeyWidth int) *keyValueTable {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(w)

	// ╭─────────┬───────────╮
	// 12--key--345--value--67
	// ╰─────────┴───────────╯
	const tablePadding = 7

	truncationWidth := desiredMaxWidth - maxKeyWidth - tablePadding
	if truncationWidth < 0 || desiredMaxWidth == 0 {
		truncationWidth = 0
	}
	return &keyValueTable{
		t:               t,
		truncationWidth: truncationWidth,
	}
}

func (k *keyValueTable) Section(rows ...row) {
	for _, e := range rows {
		if !e.Skip {
			k.t.AppendRow(table.Row{
				e.Key,
				truncate(e.Value, k.truncationWidth),
			})
		}
	}
	k.t.AppendSeparator()
}

func (k *keyValueTable) Render() {
	k.t.Render()
}
