package service

import (
	"fmt"
	"time"
)

func GetTimeWithTimezone(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %w", err)
	}

	return time.Now().In(loc), nil
}
