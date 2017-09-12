package astro

import (
	"github.com/soniakeys/meeus/deltat"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonposition"
	"github.com/soniakeys/meeus/rise"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/unit"
	"time"
)

type MoonBody struct {
	// day data
	Date        float64 // julian date at which the numbers below have been calculated
	ec          EclipticPosition
	eq          EquatorialPosition
	Events      [3]TransitEvent // rise-set-rise triplet spreading over multiple days
	Waxing      bool
	NearestFull time.Time
	NearestNew  time.Time
	// instant data
	Moment      float64 // time for which the numbers below have been calculated
	Illuminated float64 // percent
}

func NewMoon(t time.Time) *MoonBody {
	jd0 := julian.CalendarGregorianToJD(t.Year(), int(t.Month()), float64(t.Day()))
	return &MoonBody{
		Date:   jd0,
		Events: [3]TransitEvent{},
	}
}

func (moon MoonBody) ComputeEclPos(t time.Time) EclipticPosition {
	jde := julian.TimeToJD(t)
	long, lat, dist := moonposition.Position(jde)
	return EclipticPosition{
		long:     long,
		lat:      lat,
		distance: dist,
		jd:       jde,
	}
}

func (moon *MoonBody) EclipticPosition() EclipticPosition {
	return (*moon).ec
}

func (moon *MoonBody) EquatorialPosition() EquatorialPosition {
	return (*moon).eq
}

// Ecliptic position must be set before
func (moon MoonBody) ComputeEquaPos() EquatorialPosition {
	if moon.ec.jd == 0 {
		// TODO : better error handling
		return EquatorialPosition{}
	}
	α, δ := EclToEqu(&moon)
	return EquatorialPosition{
		RA:  α,
		Dec: δ,
		jd:  moon.ec.jd,
	}
}

// must set jd & equatorial position before
func (moon MoonBody) ComputeTransit(pos globe.Coord) (tRise, tTransit, tSet unit.Time, err error) {
	α := make([]unit.RA, 3)
	δ := make([]unit.Angle, 3)

	yesterday := julian.JDToTime(moon.Date).Add(-24 * time.Hour)
	moon_yesterday := NewMoon(yesterday)
	moon_yesterday.SetPositions(yesterday)

	tomorrow := julian.JDToTime(moon.Date).Add(24 * time.Hour)
	moon_tomorrow := NewMoon(tomorrow)
	moon_tomorrow.SetPositions(tomorrow)

	α[0] = moon_yesterday.eq.RA
	α[1] = moon.eq.RA
	α[2] = moon_tomorrow.eq.RA
	δ[0] = moon_yesterday.eq.Dec
	δ[1] = moon.eq.Dec
	δ[2] = moon_tomorrow.eq.Dec
	Th0 := sidereal.Apparent0UT(moon.Date)
	ΔT := deltat.PolyAfter2000(float64(julian.JDToTime(moon.Date).Year()))
	h0 := rise.Stdh0Lunar(moonposition.Parallax(moon.ec.distance))
	return rise.Times(pos, ΔT, h0, Th0, α, δ)
}

// a helper to call two methods at a time
// jd must be set before
func (moon *MoonBody) SetPositions(t time.Time) {
	if moon.Date == 0 {
		return
	}
	moon.ec = moon.ComputeEclPos(t)
	moon.eq = moon.ComputeEquaPos()
}
