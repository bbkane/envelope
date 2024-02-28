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

type kv struct {
	Key   string
	Value string
}

func tableInit(w io.Writer) table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(w)
	return t
}

func tableAddSection(t table.Writer, kvs []kv) {
	for _, e := range kvs {
		t.AppendRow(table.Row{
			e.Key,
			e.Value,
		})
	}
	t.AppendSeparator()
}
