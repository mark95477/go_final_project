package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	startDate, err := time.Parse(formatTime, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date")
	}

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, startDate.Location())
	current := startDate

	rule := strings.Split(repeat, " ")
	switch rule[0] {
	case "d":
		if len(rule) < 2 {
			return "", fmt.Errorf("missing days interval")
		}

		days, err := strconv.Atoi(rule[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid days interval")
		}

		for {
			current = current.AddDate(0, 0, days)
			if current.After(nowDate) {
				break
			}
		}

	case "y":
		for {
			current = current.AddDate(1, 0, 0)
			if current.After(nowDate) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}

	return current.Format(formatTime), nil
}
