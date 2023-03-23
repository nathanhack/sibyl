package internal

import "time"

func AtMidnight() <-chan time.Time {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	return time.After(time.Until(midnight))
}
