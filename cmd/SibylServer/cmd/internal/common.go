package internal

import "time"

func tomorrowAt6AM() time.Time {
	return time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1).Add(6 * time.Hour)
}
