package solar

import "time"

func ApproxSunsetJerusalem(t time.Time) time.Time {
	year, month, day := t.Date()
	loc := t.Location()

	hour := 18
	minute := 0

	switch month {
	case time.November, time.December, time.January:
		hour = 16
		minute = 45
	case time.February:
		hour = 17
		minute = 15
	case time.March:
		hour = 17
		minute = 45
	case time.April:
		hour = 18
		minute = 15
	case time.May:
		hour = 18
		minute = 35
	case time.June, time.July:
		hour = 18
		minute = 45
	case time.August:
		hour = 18
		minute = 25
	case time.September:
		hour = 17
		minute = 50
	case time.October:
		hour = 17
		minute = 15
	}

	return time.Date(year, month, day, hour, minute, 0, 0, loc)
}
