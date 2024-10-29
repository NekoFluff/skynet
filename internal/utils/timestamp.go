package utils

import (
	"fmt"
	"time"
)

func ConvertToUnixTimestamp(dateTimeString string, timezone string) (int64, error) {
	// Parse the date time string
	layouts := []string{"January 2, 2006 3PM", "Jan 2, 2006 3PM", "1/2/2006", time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano}
	t, err := parseTime(layouts, dateTimeString, timezone)
	if err != nil {
		return 0, err
	}

	// Return the unix timestamp
	return t.Unix(), nil
}

func parseTime(layout interface{}, value string, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	switch layout := layout.(type) {
	case string:
		t, err := time.ParseInLocation(layout, value, loc)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	case []string:
		for _, l := range layout {
			t, err := time.ParseInLocation(l, value, loc)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("could not parse time with any of the provided layouts")
	default:
		return time.Time{}, fmt.Errorf("invalid layout type")
	}
}
