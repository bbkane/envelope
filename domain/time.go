package domain

import "time"

// -- Utility function

// TimePtrToStringPtr uses TimeToString to optonally format a string
func TimePtrToStringPtr(t *time.Time) *string {
	var s *string
	if t != nil {
		tmp := TimeToString(*t)
		s = &tmp
	}
	return s
}

// TimeToString converts a time to UTC, then formats as RFC3339
func TimeToString(t time.Time) string {

	return t.Round(0).UTC().Format(time.RFC3339)
}

// StringToTime converts a RFC3339 formatted string into a time.Time
func StringToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return t, err
	}
	return t.Round(0), nil
}

// StringToTimeMust works like StringToTime but panics on errors.
// I think this is usually acceptable as times are formatted pretty carefully
// in the db
func StringToTimeMust(s string) time.Time {
	t, err := StringToTime(s)
	if err != nil {
		panic(err)
	}
	return t.Round(0)
}
