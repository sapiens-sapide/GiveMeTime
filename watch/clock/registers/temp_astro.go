package registers

import (
	"github.com/soniakeys/unit"
	"time"
)

var (
	SunRiseTime      time.Time
	SunRiseAz        unit.Angle
	SunZenithTime    time.Time
	SunZenithAlt     unit.Angle
	SunSetTime       time.Time
	SunSetAz         unit.Angle
	SunlightDuration time.Duration
	SunAz            unit.Angle
	MoonPercent      float64
	MoonRise         time.Time
	MoonSet          time.Time
	MoonAz           unit.Angle
)
