package tableprint

func valOrEmpty(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
