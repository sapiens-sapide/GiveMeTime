// This package build ephemeris data for a given number of days and a given geo-position
package ephemeris

import (
	"fmt"
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	pp "github.com/soniakeys/meeus/planetposition"
	"os"
	"time"
)

// all data we need for one day
// all times are given in minutes, starting at midnight
type DayEphemeris struct {
	Date              time.Time
	Sun               SunEphemeris
	Moon              MoonEphemeris
	EquinoxOrSolstice bool // true if solstice or equinox is within the current day
}

// sun ephemeris for one day
// all times are given in minutes, starting at midnight
type SunEphemeris struct {
	Rise      float32
	RiseAz    float32 // in degrees
	CivilRise float32
	Set       float32
	SetAz     float32 // in degrees
	CivilSet  float32
	Zenith    float32 // noon time
}

// moon ephemeris for one day
type MoonEphemeris struct {
}

func EphemerisForDay(day time.Time, pos globe.Coord) (eph DayEphemeris, err error) {
	os.Setenv("VSOP87", "/usr/local/goland/src/github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro")
	astro.Earth, err = pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}
	astro.Position = &pos
	// sun computations
	astro.Sun = astro.NewSun(day)
	astro.Sun.SetPositions(julian.JDToTime(astro.Sun.Date).Add(3 * time.Hour)) // compute position at 3 o'clock in the morning to be closer to rise/transit/set events ?

	utcRise, utcTransit, utcSet, err := astro.Sun.ComputeTransit(*astro.Position)
	if err != nil {
		fmt.Printf("Error when computing sun's transit data : %s\n", err)
		return
	}
	utcRiseCiv, _, utcSetCiv, err := astro.Sun.ComputeTransit(*astro.Position, 350)
	_, tz := day.Zone()
	eph.Date = day
	eph.Sun.Rise = float32(utcRise.Sec()+float64(tz)) / 60.0
	eph.Sun.Zenith = float32(utcTransit.Sec()+float64(tz)) / 60.0
	eph.Sun.Set = float32(utcSet.Sec()+float64(tz)) / 60.0
	eph.Sun.CivilRise = float32(utcRiseCiv.Sec()+float64(tz)) / 60.0
	eph.Sun.CivilSet = float32(utcSetCiv.Sec()+float64(tz)) / 60.0
	return
}
