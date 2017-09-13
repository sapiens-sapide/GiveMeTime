package computers

import (
	"fmt"
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
	r "github.com/sapiens-sapide/GiveMeTime/watch/clock/registers"
	"github.com/soniakeys/meeus/base"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonillum"
	pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/unit"
	"os"
	"time"
)

type nextDay struct{}

var NextDay = new(nextDay)

func (nd nextDay) Trigger() {
	ny, nm, nday, _ := NextDate(int(r.Year.Status().(uint16)), int(r.Month.Status().(uint8)), int(r.Day.Status().(uint16)))
	r.SetDate(ny, nm, nday)
	DayComputation(r.Now())
}

func (nd nextDay) Status() interface{} { return nil }

func (nd nextDay) Set(value interface{}) {}

// returns next logical date (day, weekday, month, year) following the given date (ie the "tomorrow")
func NextDate(y, m, d int) (ny, nm, nd, nwd int) {
	next := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	ny = next.Year()
	nm = int(next.Month())
	nd = next.Day()
	wd := int(next.Weekday())
	if wd == 0 {
		nwd = 7
	} else {
		nwd = wd
	}
	return
}

func DayComputation(now time.Time) {
	os.Setenv("VSOP87", "/usr/local/goland/src/github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro")
	var err error
	astro.Earth, err = pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*astro.Position = &globe.Coord{
		Lat: unit.NewAngle('+', 43, 29, 11.0),
		Lon: unit.NewAngle('+', 110, 45, 22.9392), //positive westward
	}*/

	astro.Position = &globe.Coord{
		Lat: unit.NewAngle('+', 48, 51, 39),
		Lon: unit.NewAngle('-', 2, 22, 1), //positive westward
	}

	// sun computations
	astro.Sun = astro.NewSun(now)
	astro.Sun.SetPositions(julian.JDToTime(astro.Sun.Date).Add(3 * time.Hour)) // compute position at 3 o'clock in the morning to be closer to rise/transit/set events ?

	// TODO : these times are UTC
	tRise, tTransit, tSet, err := astro.Sun.ComputeTransit(*astro.Position)
	if err != nil {
		fmt.Printf("Error when computing sun's transit data : %s\n", err)
		return
	}

	r.SunRiseTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tRise)) * time.Second).Add(time.Duration(r.Tz.Status().(int8)) * time.Hour)
	r.SunZenithTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tTransit)) * time.Second).Add(time.Duration(r.Tz.Status().(int8))  * time.Hour)
	r.SunSetTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tSet)) * time.Second).Add(time.Duration(r.Tz.Status().(int8))  * time.Hour)

	sunrise_jd := julian.TimeToJD(r.SunRiseTime)
	sunrise_az, _ := coord.EqToHz(astro.Sun.EquatorialPosition().RA, astro.Sun.EquatorialPosition().Dec, astro.Position.Lat, astro.Position.Lon, sidereal.Apparent(sunrise_jd))
	r.SunRiseAz = unit.AngleFromDeg(sunrise_az.Deg() + 180)

	sunzenith_jd := julian.TimeToJD(r.SunZenithTime)
	_, r.SunZenithAlt = coord.EqToHz(astro.Sun.EquatorialPosition().RA, astro.Sun.EquatorialPosition().Dec, astro.Position.Lat, astro.Position.Lon, sidereal.Apparent(sunzenith_jd))

	sunset_jd := julian.TimeToJD(r.SunSetTime)
	sunset_az, _ := coord.EqToHz(astro.Sun.EquatorialPosition().RA, astro.Sun.EquatorialPosition().Dec, astro.Position.Lat, astro.Position.Lon, sidereal.Apparent(sunset_jd))
	r.SunSetAz = unit.AngleFromDeg(sunset_az.Deg() + 180)
	r.SunlightDuration = r.SunSetTime.Sub(r.SunRiseTime)

	// moon computations
	astro.Moon = astro.NewMoon(now)

	/* TODO
	fmt.Println(julian.JDToCalendar(moonphase.Full(base.JDEToJulianYear(jd0-0.125)))) //plus proche pleine lune
	fmt.Println(julian.JDToCalendar(moonphase.New(base.JDEToJulianYear(jd0-0.125)))) //plus proche nouvelle lune
	*/
	astro.Moon.SetPositions(julian.JDToTime(astro.Moon.Date))

	// TODO : these times are UTC
	tRise, tTransit, tSet, err = astro.Moon.ComputeTransit(*astro.Position)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.MoonRise = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tRise)) * time.Second).Add(time.Duration(r.Tz.Status().(int8))  * time.Hour)
	r.MoonSet = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tSet)) * time.Second).Add(time.Duration(r.Tz.Status().(int8))  * time.Hour)

	i := moonillum.PhaseAngle3(astro.Moon.Date)
	r.MoonPercent = base.Illuminated(i)
}
