package utils

import (
	"fmt"
	"log"
	"time"
)

func FormatDateToIndonesianFormat(t time.Time) (string, error) {

	utcTime := t.UTC()

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading timezone: %v", err)
		return "", fmt.Errorf("could not load location 'Asia/Jakarta'")
	}

	indonesianTime := utcTime.In(loc)

	return indonesianTime.Format("02-01-2006 15:04"), nil
}