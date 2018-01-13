package astro

import (
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonphase"
	"math"
	"time"
)

// IsMoonEvent returns
// 0 if no special moon event for given date
// 1 if new moon (± 24 hours)
// 2 if full moon (± 24 hours)
func IsMoonEvent(obs Observer) uint8 {
	year := float64(obs.Date.Year()) + (float64(obs.Date.YearDay()) / 365.0)
	fullmoon := julian.JDToTime(moonphase.Full(year))
	newmoon := julian.JDToTime(moonphase.New(year))
	if math.Abs(float64(fullmoon.Sub(obs.Date)/time.Hour)) < 24 {
		return 2
	}
	if math.Abs(float64(newmoon.Sub(obs.Date)/time.Hour)) < 24 {
		return 1
	}
	return 0
}
