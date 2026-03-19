package biblical

import "time"

func IsAfterSunset(now, sunset time.Time) bool {
	return now.Equal(sunset) || now.After(sunset)
}

func GetBiblicalDate(now, sunset time.Time) time.Time {
	if IsAfterSunset(now, sunset) {
		return now.AddDate(0, 0, 1)
	}
	return now
}
