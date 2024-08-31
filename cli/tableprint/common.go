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

func tableInit(w io.Writer) table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(w)
	return t
}

// tableAddSection adds a section to the table with the given key-value pairs and then a separator. If a value is empty, the row is not added.
func tableAddSection(t table.Writer, rows []row) {
	for _, e := range rows {
		if !e.Skip {
			t.AppendRow(table.Row{
				e.Key,
				e.Value,
			})
		}
	}
	t.AppendSeparator()
}
