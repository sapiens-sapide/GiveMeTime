package registers

import (
	"github.com/soniakeys/unit"
	"time"
	"github.com/soniakeys/meeus/globe"
)

var (
	SunRiseTime      time.Time
	SunRiseAz        unit.Angle
	SunZenithTime    time.Time
	SunZenithAlt     unit.Angle
	SunSetTime       time.Time
	SunSetAz         unit.Angle
	SunlightDuration time.Duration
	Position         globe.Coord
	SunAz            unit.Angle
)
