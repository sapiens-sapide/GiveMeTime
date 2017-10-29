// This package build ephemeris data for a given number of days and a given geo-position
package ephemeris

import (
	"fmt"
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
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
// all times are given in seconds, starting at midnight
type SunEphemeris struct {
	Rise      uint32
	RiseAz    float32 // in degrees
	CivilRise uint32
	Set       uint32
	SetAz     float32 // in degrees
	CivilSet  uint32
	Zenith    uint32 // noon time
}

// moon ephemeris for one day
type MoonEphemeris struct {
}

func EphemerisForDay(t time.Time, lat, lon float64) (eph DayEphemeris, err error) {
	eph.Date = t.Add(-(time.Duration(t.Hour()) * time.Hour) - (time.Duration(t.Minute()) * time.Minute) - (time.Duration(t.Second()) * time.Second))
	o := astro.NewObserver(eph.Date, lat, lon)
	if o == nil {
		fmt.Println("Error when creating observer")
		return
	}
	sun := astro.NewSun(o)
	rise, transit, set, err := sun.ComputeTransit(astro.SunStdAlt)
	if err != nil {
		fmt.Printf("Error when computing sun's transit data : %s\n", err)
		return
	}
	eph.Sun.Rise = uint32(rise)
	eph.Sun.Zenith = uint32(transit)
	eph.Sun.Set = uint32(set)

	// compute apparent sun at sunrise time
	o_rise := astro.NewObserver(eph.Date.Add(time.Duration(uint64(rise))*time.Second), 48.860833, -2.366944)
	sun_rise := astro.NewSun(o_rise)
	apparent_rise := sun_rise.ApparentPosition()
	eph.Sun.RiseAz = float32(apparent_rise.Az)

	// compute apparent sun at sunset time
	o_set := astro.NewObserver(eph.Date.Add(time.Duration(uint64(set))*time.Second), 48.860833, -2.366944)
	sun_set := astro.NewSun(o_set)
	apparent_set := sun_set.ApparentPosition()
	eph.Sun.SetAz = float32(apparent_set.Az)

	rise, _, set, err = sun.ComputeTransit(astro.SunCivilAlt)
	if err != nil {
		fmt.Printf("Error when computing sun's transit data : %s\n", err)
		return
	}
	eph.Sun.CivilRise = uint32(rise)
	eph.Sun.CivilSet = uint32(set)

	return
}
