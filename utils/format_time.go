package utils

import (
	"log"
	"time"
)

func FormatDateToIndonesianFormat(t time.Time) string {

	utcTime := t.UTC()

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading timezone: %v", err)
		return ""
	}

	indonesianTime := utcTime.In(loc)

	return indonesianTime.Format("02-01-2006 15:04")
}
