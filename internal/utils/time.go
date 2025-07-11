package utils

import (
	"fmt"
	"time"
)

func RelativeTime(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < 0 {
		return "in the future"
	}

	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := hours / 24
	weeks := days / 7
	months := days / 30
	years := days / 365

	switch {
	case years > 0:
		return fmt.Sprintf("%d year%s ago", years, plural(years))
	case months > 0:
		return fmt.Sprintf("%d month%s ago", months, plural(months))
	case weeks > 0:
		return fmt.Sprintf("%d week%s ago", weeks, plural(weeks))
	case days > 0:
		return fmt.Sprintf("%d day%s ago", days, plural(days))
	case hours > 0:
		return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
	case minutes > 0:
		return fmt.Sprintf("%d minute%s ago", minutes, plural(minutes))
	default:
		return fmt.Sprintf("%d second%s ago", seconds, plural(seconds))
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
