package utils

import (
	"fmt"
	"time"
)

func ConvertToUnixTimestamp(input string) (int64, error) {
	// Parse the date time string
	layouts := []string{"January 2, 2006 3PM", "Jan 2, 2006 3PM", "01/02/2006", time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano}
	t, err := parseTime(layouts, input)
	if err != nil {
		return 0, err
	}

	// Return the unix timestamp
	return t.Unix(), nil
}

func parseTime(layout interface{}, value string) (time.Time, error) {
	switch layout := layout.(type) {
	case string:
		return time.Parse(layout, value)
	case []string:
		for _, l := range layout {
			t, err := time.Parse(l, value)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("could not parse time with any of the provided layouts")
	default:
		return time.Time{}, fmt.Errorf("invalid layout type")
	}
}
