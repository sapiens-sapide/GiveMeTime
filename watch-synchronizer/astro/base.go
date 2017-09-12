package astro

import (
	"time"
	"math"
)


func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// returns a time.Time from a decimal day
// current time.Date is set by default from time.Now() if none is provided.
func TimeFromDay(d float64, t... time.Time) time.Time {
	seconds := int64(d * 24 * 60 * 60)
	var local_time time.Time
	if len(t) != 1 {
		local_time = time.Now()
	} else {
		local_time = t[0]
	}

	return time.Date(local_time.Year(), local_time.Month(), local_time.Day(), 0, 0, 0, 0, local_time.Location()).Add(time.Duration(seconds) * time.Second)
}

// returns a decimal day from a time
func DayFromTime(t time.Time) float64 {
	return float64(t.Hour()*3600 + t.Minute()*60 + t.Second())/86400.0
}