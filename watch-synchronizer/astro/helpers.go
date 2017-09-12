package astro

import (
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/nutation"
	"github.com/soniakeys/unit"
	"math"
	"time"
)

func DegToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func RadToDeg(radians float64) float64 {
	return radians * 180 / math.Pi
}

func RadToHour(radians float64) float64 {
	return radians * (180 / (15 * math.Pi))
}

func RadToDay(radians float64) float64 {
	return RadToHour(radians) / 24
}

// returns a time.Time from a decimal day
// current time.Date is set by default from time.Now() if none is provided.
func TimeFromDay(d float64, t ...time.Time) time.Time {
	seconds := int64(d * 24 * 60 * 60)
	var local_time time.Time
	if len(t) != 1 {
		local_time = time.Now()
	} else {
		local_time = t[0]
	}

	return time.Date(local_time.Year(), local_time.Month(), local_time.Day(), 0, 0, 0, 0, local_time.Location()).Add(time.Duration(seconds) * time.Second)
}

func EclToEqu(cb CelestialBody) (α unit.RA, δ unit.Angle) {
	ecl := cb.EclipticPosition()
	Δψ, Δε := nutation.Nutation(ecl.jd)
	a := unit.AngleFromSec(-20.4898).Div(ecl.distance)
	λ := ecl.long + Δψ + a
	ε := nutation.MeanObliquityLaskar(ecl.jd) + Δε
	sε, cε := ε.Sincos()
	return coord.EclToEq(λ, ecl.lat, sε, cε)
}

// returns a decimal day from a time
func DayFromTime(t time.Time) float64 {
	return float64(t.Hour()*3600+t.Minute()*60+t.Second()) / 86400.0
}
