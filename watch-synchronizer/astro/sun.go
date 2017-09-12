package astro

import (
	"fmt"
	"github.com/soniakeys/meeus/deltat"
	"github.com/soniakeys/meeus/eqtime"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/meeus/rise"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/meeus/solar"
	"github.com/soniakeys/unit"
	"math"
	"os"
	"time"
)

type SunBody struct {
	// day data
	Date   float64 // julian date at which the numbers below have been calculated
	ec     EclipticPosition
	eq     EquatorialPosition
	Events [3]TransitEvent // rise, zenith, set
	// instant data
	Moment         float64 // time for which the numbers below have been calculated
	Az, Elev, Dist float64
}

func NewSun(t time.Time) *SunBody {
	// set time at beginning of the day for astronomical calculation.
	jd0 := julian.CalendarGregorianToJD(t.Year(), int(t.Month()), float64(t.Day()))
	return &SunBody{
		Date:   jd0,
		Events: [3]TransitEvent{},
	}
}

func (sun SunBody) ComputeEclPos(t time.Time) EclipticPosition {
	jde := julian.TimeToJD(t)
	long, lat, dist := solar.TrueVSOP87(Earth, jde)
	return EclipticPosition{
		long:     long,
		lat:      lat,
		distance: dist,
		jd:       jde,
	}
}

// Ecliptic position must be set before
func (sun SunBody) ComputeEquaPos() EquatorialPosition {
	if sun.ec.jd == 0 {
		// TODO : better error handling
		return EquatorialPosition{}
	}
	α, δ := EclToEqu(&sun)
	return EquatorialPosition{
		RA:  α,
		Dec: δ,
		jd:  sun.ec.jd,
	}
}

// must set jd & equatorial position before
func (sun SunBody) ComputeTransit(pos globe.Coord) (tRise, tTransit, tSet unit.Time, err error) {
	α := make([]unit.RA, 3)
	δ := make([]unit.Angle, 3)

	yesterday := julian.JDToTime(sun.Date).Add(-24 * time.Hour)
	sun_yesterday := NewSun(yesterday)
	sun_yesterday.SetPositions(yesterday)

	tomorrow := julian.JDToTime(sun.Date).Add(24 * time.Hour)
	sun_tomorrow := NewSun(tomorrow)
	sun_tomorrow.SetPositions(tomorrow)

	α[0] = sun_yesterday.eq.RA
	α[1] = sun.eq.RA
	α[2] = sun_tomorrow.eq.RA
	δ[0] = sun_yesterday.eq.Dec
	δ[1] = sun.eq.Dec
	δ[2] = sun_tomorrow.eq.Dec
	Th0 := sidereal.Apparent0UT(sun.Date)
	ΔT := deltat.PolyAfter2000(float64(julian.JDToTime(sun.Date).Year()))
	h0 := rise.Stdh0Solar
	return rise.Times(pos, ΔT, h0, Th0, α, δ)
}

// a helper to call two methods at a time
// jd must be set before
func (sun *SunBody) SetPositions(t time.Time) {
	if sun.Date == 0 {
		return
	}
	sun.ec = sun.ComputeEclPos(t)
	sun.eq = sun.ComputeEquaPos()
}

func (sun *SunBody) EclipticPosition() EclipticPosition {
	return (*sun).ec
}

func (sun *SunBody) EquatorialPosition() EquatorialPosition {
	return (*sun).eq
}

// compute data for given day
// and set sun's properties with results of these computations.
func (sun *SunBody) SetDay(t time.Time) {
	sun.ComputeEclPos(t)
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
