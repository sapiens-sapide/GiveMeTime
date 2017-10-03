package astro

import (
	"fmt"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/eqtime"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/meeus/solar"
	"math"
	"os"
	"time"
	"github.com/soniakeys/unit"
	"github.com/soniakeys/meeus/deltat"
	"github.com/soniakeys/meeus/rise"
)

type SunBody struct {
	observer *Observer
	ec       EclipticPosition
	eq       EquatorialPosition
	ap       ApparentPosition
}

// NewSun returns a SunBody for given observer
func NewSun(o *Observer) *SunBody {
	s := new(SunBody)
	s.observer = o
	s.ComputePositions()
	return s
}

// a helper to call tree methods at a time
// positions must be computed in order
func (sun *SunBody) ComputePositions() {
	sun.ec = sun.ComputeEclPos()
	sun.eq = sun.ComputeEquaPos()
	sun.ap = sun.ComputeApparentPos()
}

func (sun SunBody) ComputeEclPos() EclipticPosition {
	long, lat, dist := solar.TrueVSOP87(sun.observer.earth, sun.observer.JulianDate)
	return EclipticPosition{
		long:     long,
		lat:      lat,
		distance: dist,
	}
}

// Ecliptic position must be computed before
func (sun SunBody) ComputeEquaPos() EquatorialPosition {
	α, δ := EclToEqu(&sun)
	return EquatorialPosition{
		RA:  α,
		Dec: δ,
	}
}

// Equatorial position must be computed before
func (sun SunBody) ComputeApparentPos() ApparentPosition {
	az, el := coord.EqToHz(sun.eq.RA, sun.eq.Dec, sun.observer.Position.Lat, sun.observer.Position.Lon, sidereal.Apparent(sun.JD()))
	return ApparentPosition{
		Az: az.Deg() + 180,
		Elev: el.Deg(),
	}
}

// must set jd & equatorial position before
// returned times are in seconds from midnight
// alt is to adjust twilight rise/set. It must be given in minutes of angle (there are 60 min in a degree).
func (sun SunBody) ComputeTransit(alt float64) (tRise, tTransit, tSet unit.Time, err error) {
	α := make([]unit.RA, 3)
	δ := make([]unit.Angle, 3)

	yesterday := julian.JDToTime(sun.JD()).Add(-24 * time.Hour)
	o_y := NewObserver(yesterday, sun.observer.Position.Lat.Deg(), sun.observer.Position.Lon.Deg())
	sun_yesterday := NewSun(o_y)

	tomorrow := julian.JDToTime(sun.JD()).Add(24 * time.Hour)
	o_t := NewObserver(tomorrow, sun.observer.Position.Lat.Deg(), sun.observer.Position.Lon.Deg())
	sun_tomorrow := NewSun(o_t)

	α[0] = sun_yesterday.eq.RA
	α[1] = sun.eq.RA
	α[2] = sun_tomorrow.eq.RA
	δ[0] = sun_yesterday.eq.Dec
	δ[1] = sun.eq.Dec
	δ[2] = sun_tomorrow.eq.Dec
	Th0 := sidereal.Apparent0UT(sun.JD())
	ΔT := deltat.PolyAfter2000(float64(julian.JDToTime(sun.JD()).Year()))
	h0 := unit.AngleFromMin(alt)
	return rise.Times(*sun.observer.Position, ΔT, h0, Th0, α, δ) //TODO : improve error handling, notably when error is 'Circumpolar'
}

// CelestialBody interface
func (sun *SunBody) JD() float64 {
	return (*sun).observer.JulianDate
}

func (sun *SunBody) EclipticPosition() EclipticPosition {
	return (*sun).ec
}

func (sun *SunBody) EquatorialPosition() EquatorialPosition {
	return (*sun).eq
}

func (sun *SunBody) ApparentPosition() ApparentPosition {
	return (*sun).ap
}

// returns solar hour angle at sunrise
// lat is observer latitude (positive northward)
// sun_dec is declination of the sun for the day
// algorithm from NOAA spreadsheet (https://www.esrl.noaa.gov/gmd/grad/solcalc/NOAA_Solar_Calculations_day.xls)
func SunRiseHA(lat, sun_dec float64) float64 {
	return RadToDeg(math.Acos(
		math.Cos(DegToRad(90.833))/
			(math.Cos(DegToRad(lat))*math.Cos(DegToRad(sun_dec))) -
			math.Tan(DegToRad(lat))*math.Tan(DegToRad(sun_dec))))
}

func SolarNoon(day time.Time, long, eot float64, tz int) time.Time {
	noon := (720 - 4*long - eot + float64(tz)*60) / 1440
	second_shift := noon * 24 * 60 * 60
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Add(time.Duration(second_shift) * time.Second)
}

func TrueSolarTime(t time.Time, pos globe.Coord) time.Time {
	os.Setenv("VSOP87", "./astro")
	jd := julian.TimeToJD(t)
	earth, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	eq := eqtime.E(jd, earth)
	_, offset := t.Zone()
	tst := math.Mod(float64(t.Hour())*60.0+float64(t.Minute())+(float64(t.Second())/60.0)+eq.Min()+(4*(-pos.Lon.Deg())-60.0*float64(offset)/3600.0), 1440)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(time.Duration(tst*60) * time.Second)
}

// tst = True Solar Time
func SolarHourAngle(tst time.Time) float64 {
	tst_minutes := float64(tst.Hour()*60+tst.Minute()) + float64(tst.Second())/60.0
	t := tst_minutes / 4.0
	if t < 0 {
		return t + 180
	}
	return t - 180
}
