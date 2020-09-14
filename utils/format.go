package utils

import (
	"fmt"
	"math"
	"time"
)

const (
	millisInSecond  = 1000
	secondsInMinute = 60
	minutesInHour   = 60
)

// FormatDuration will formart time output
//
// If the time duration is no longer than 10s, it will print decimal if showDecimal is true.
// If the time duration is no longer than 1min, it will print the seconds.
// If the time duration is no longer than 1h, it will print as min:sec (each with two width).
// If the time duration is longer than 1h, it will print like h:m:s (each with two width).
func FormatDuration(duration time.Duration, showDecimals bool) string {
	if duration < 10*time.Second && showDecimals {
		return fmt.Sprintf("%.1fs", float64(duration.Milliseconds())/millisInSecond)
	}
	seconds := int(math.Floor(duration.Seconds())) % secondsInMinute
	if duration < time.Minute {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := int(math.Floor(duration.Minutes())) % minutesInHour
	if duration < time.Hour {
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	hours := int(math.Floor(duration.Hours()))
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
