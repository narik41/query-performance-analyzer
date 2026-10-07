package util

import "time"

// TimeNow returns the current time
func TimeNow() time.Time {
	return time.Now()
}

// TimeSince returns milliseconds since the given time
func TimeSince(start time.Time) float64 {
	return float64(time.Since(start).Microseconds()) / 1000.0
}