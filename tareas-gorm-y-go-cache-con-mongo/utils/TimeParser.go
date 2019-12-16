package utils

import "time"

func TimeMustParse(dateStr string) time.Time {
	miTime, _ := time.Parse("2006-01-02", dateStr)
	return miTime
}
